// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// SPDX-License-Identifier: AGPL-3.0-only

// via54Design — Image Analysis Engine
// Pure Go image analysis: color, brightness, contrast, detail, mood/style suggestion.
// Replaces scripts/img2prompt.py — zero external dependencies.
package vision

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"sort"
)

// ─── Types ────────────────────────────────────────────────────────────────

// ColorEntry represents a quantized dominant color.
type ColorEntry struct {
	RGB   []int  `json:"rgb"`
	Hex   string `json:"hex"`
	Count int    `json:"count"`
}

// ImageAnalysis holds all extracted visual features from an image.
type ImageAnalysis struct {
	Width           int          `json:"width"`
	Height          int          `json:"height"`
	AspectRatio     string       `json:"aspect_ratio"`
	Orientation     string       `json:"orientation"`
	Brightness      float64      `json:"brightness"`
	BrightnessLabel string       `json:"brightness_label"`
	Contrast        string       `json:"contrast"`
	Colorfulness    string       `json:"colorfulness"`
	Detail          string       `json:"detail"`
	DominantColors  []ColorEntry `json:"dominant_colors"`
	SuggestedMoods  []string     `json:"suggested_moods"`
	SuggestedStyles []string     `json:"suggested_styles"`
	GeneratedPrompt string       `json:"generated_prompt,omitempty"`
	Error           string       `json:"error,omitempty"`
}

// ─── Public API ───────────────────────────────────────────────────────────

// AnalyzeImage extracts visual features from an image at the given path.
// Supports PNG, JPEG, GIF. Returns a JSON-serializable struct.
func AnalyzeImage(path string) *ImageAnalysis {
	f, err := os.Open(path)
	if err != nil {
		return &ImageAnalysis{Error: fmt.Sprintf("file not found: %s", path)}
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return &ImageAnalysis{Error: fmt.Sprintf("decode failed: %v", err)}
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w == 0 || h == 0 {
		return &ImageAnalysis{Error: "image has zero dimensions"}
	}

	// Aspect ratio
	aspect := float64(w) / float64(h)
	orientation := "square"
	if aspect > 1.1 {
		orientation = "landscape"
	} else if aspect < 0.9 {
		orientation = "portrait"
	}

	// Resize to 64x64 for color analysis
	small := resizeNN(img, 64, 64)
	pixels := samplePixels(small)

	// Dominant colors (quantize to 16-colors, count)
	quantized := quantizeColors(pixels)
	topColors := topColorEntries(quantized, 5)

	// Grayscale + brightness
	grayPixels := make([]float64, len(pixels))
	var avgBrightness float64
	for i, c := range pixels {
		g := luminance(float64(c.R), float64(c.G), float64(c.B))
		grayPixels[i] = g
		avgBrightness += g
	}
	avgBrightness /= float64(len(pixels))

	brightnessLabel := "mid-key"
	if avgBrightness > 200 {
		brightnessLabel = "high-key"
	} else if avgBrightness < 60 {
		brightnessLabel = "low-key"
	}

	// Contrast (std dev of grayscale)
	var variance float64
	for _, g := range grayPixels {
		d := g - avgBrightness
		variance += d * d
	}
	if len(grayPixels) > 0 {
		variance /= float64(len(grayPixels))
	}
	contrastStd := math.Sqrt(variance)

	contrast := "moderate"
	if contrastStd > 60 {
		contrast = "high contrast"
	} else if contrastStd < 25 {
		contrast = "soft"
	}

	// Colorfulness: variance across RGB channels
	rPixels := make([]float64, len(pixels))
	gPixels := make([]float64, len(pixels))
	bPixels := make([]float64, len(pixels))
	for i, c := range pixels {
		rPixels[i] = float64(c.R)
		gPixels[i] = float64(c.G)
		bPixels[i] = float64(c.B)
	}
	rStd := stdDev(rPixels)
	gStd := stdDev(gPixels)
	bStd := stdDev(bPixels)
	colorVar := rStd + gStd + bStd

	colorful := "balanced"
	if colorVar > 120 {
		colorful = "vibrant"
	} else if colorVar < 50 {
		colorful = "muted"
	}

	// Edge detection for detail level
	edgeDensity := computeEdgeDensity(grayPixels, 64, 64)

	detail := "moderate detail"
	if edgeDensity > 40 {
		detail = "high detail"
	} else if edgeDensity < 10 {
		detail = "smooth"
	}

	moods := suggestMoods(avgBrightness, contrastStd, colorVar)
	styles := suggestStyles(orientation, brightnessLabel, colorful)

	return &ImageAnalysis{
		Width:           w,
		Height:          h,
		AspectRatio:     fmt.Sprintf("%.2f", aspect),
		Orientation:     orientation,
		Brightness:      math.Round(avgBrightness*10) / 10,
		BrightnessLabel: brightnessLabel,
		Contrast:        contrast,
		Colorfulness:    colorful,
		Detail:          detail,
		DominantColors:  topColors,
		SuggestedMoods:  moods,
		SuggestedStyles: styles,
	}
}

// AnalyzeImageToMap wraps AnalyzeImage and returns a map[string]interface{}
// suitable for JSON serialization, matching the old Python output format.
func AnalyzeImageToMap(path string) map[string]interface{} {
	a := AnalyzeImage(path)
	if a.Error != "" {
		return map[string]interface{}{"error": a.Error}
	}
	data, _ := json.Marshal(a)
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	return m
}

// BuildPromptFromAnalysis builds a structured prompt string from analysis + user input.
func BuildPromptFromAnalysis(analysis *ImageAnalysis, userPrompt string) string {
	parts := []string{}
	if userPrompt != "" {
		parts = append(parts, userPrompt)
	}

	// Visual context
	ctx := []string{}
	if analysis.BrightnessLabel != "" {
		ctx = append(ctx, fmt.Sprintf("%s lighting", analysis.BrightnessLabel))
	}
	if analysis.Contrast != "" {
		ctx = append(ctx, analysis.Contrast)
	}
	if analysis.Colorfulness != "" {
		ctx = append(ctx, fmt.Sprintf("%s colors", analysis.Colorfulness))
	}
	if analysis.Detail != "" {
		ctx = append(ctx, analysis.Detail)
	}
	if len(ctx) > 0 {
		parts = append(parts, fmt.Sprintf("(%s)", joinComma(ctx)))
	}

	// Dominant colors as mood palette
	if len(analysis.DominantColors) > 0 {
		hexes := make([]string, 0, 4)
		for i, c := range analysis.DominantColors {
			if i >= 4 {
				break
			}
			hexes = append(hexes, c.Hex)
		}
		if len(hexes) > 0 {
			parts = append(parts, fmt.Sprintf("color palette: %s", joinComma(hexes)))
		}
	}

	// Style hints
	if len(analysis.SuggestedStyles) > 0 {
		parts = append(parts, fmt.Sprintf("--style %s", joinComma(analysis.SuggestedStyles)))
	}

	if len(parts) == 0 {
		return "a photograph with balanced composition and natural lighting"
	}

	return joinDot(parts) + "."
}

// BuildPromptFromAnalysisMap wraps BuildPromptFromAnalysis to work with map input.
func BuildPromptFromAnalysisMap(analysis map[string]interface{}, userPrompt string) string {
	// Build an ImageAnalysis from map for clean prompt generation
	a := &ImageAnalysis{}
	if v, ok := analysis["brightness_label"].(string); ok {
		a.BrightnessLabel = v
	}
	if v, ok := analysis["contrast"].(string); ok {
		a.Contrast = v
	}
	if v, ok := analysis["colorfulness"].(string); ok {
		a.Colorfulness = v
	}
	if v, ok := analysis["detail"].(string); ok {
		a.Detail = v
	}
	if v, ok := analysis["suggested_styles"].([]interface{}); ok {
		for _, s := range v {
			if str, ok := s.(string); ok {
				a.SuggestedStyles = append(a.SuggestedStyles, str)
			}
		}
	}
	if v, ok := analysis["dominant_colors"].([]interface{}); ok {
		for _, ci := range v {
			if cm, ok := ci.(map[string]interface{}); ok {
				ce := ColorEntry{}
				if h, ok := cm["hex"].(string); ok {
					ce.Hex = h
				}
				a.DominantColors = append(a.DominantColors, ce)
			}
		}
	}
	return BuildPromptFromAnalysis(a, userPrompt)
}

// ─── Helpers ──────────────────────────────────────────────────────────────

// luminance converts RGB (0-255 each) to grayscale using standard weights.
func luminance(r, g, b float64) float64 {
	return 0.299*r + 0.587*g + 0.114*b
}

// nearest-neighbor resize
func resizeNN(img image.Image, newW, newH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			sx := x * srcW / newW
			sy := y * srcH / newH
			dst.Set(x, y, img.At(sx+srcBounds.Min.X, sy+srcBounds.Min.Y))
		}
	}
	return dst
}

// samplePixels extracts all RGBA pixel values from an image into a flat slice.
func samplePixels(img *image.RGBA) []color.RGBA {
	b := img.Bounds()
	pixels := make([]color.RGBA, 0, b.Dx()*b.Dy())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: 255,
			})
		}
	}
	return pixels
}

// quantizeColors quantizes RGB values to 16-color equivalents (32-step buckets).
func quantizeColors(pixels []color.RGBA) map[color.RGBA]int {
	counts := make(map[color.RGBA]int)
	for _, p := range pixels {
		q := color.RGBA{
			R: (p.R / 32) * 32,
			G: (p.G / 32) * 32,
			B: (p.B / 32) * 32,
			A: 255,
		}
		counts[q]++
	}
	return counts
}

// topColorEntries returns the top N color entries sorted by count descending.
func topColorEntries(quantized map[color.RGBA]int, n int) []ColorEntry {
	type kv struct {
		c color.RGBA
		n int
	}
	var sorted []kv
	for c, count := range quantized {
		if count > 5 {
			sorted = append(sorted, kv{c, count})
		}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].n > sorted[j].n
	})

	if n > len(sorted) {
		n = len(sorted)
	}

	entries := make([]ColorEntry, 0, n)
	for i := 0; i < n; i++ {
		c := sorted[i].c
		entries = append(entries, ColorEntry{
			RGB:   []int{int(c.R), int(c.G), int(c.B)},
			Hex:   fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B),
			Count: sorted[i].n,
		})
	}
	return entries
}

// stdDev computes population standard deviation of a float64 slice.
func stdDev(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))
	var variance float64
	for _, v := range data {
		d := v - mean
		variance += d * d
	}
	variance /= float64(len(data))
	return math.Sqrt(variance)
}

// computeEdgeDensity computes an edge density metric using a Laplacian-like
// approximation on the grayscale 64x64 image. Higher values = more detail.
func computeEdgeDensity(gray []float64, w, h int) float64 {
	if len(gray) < 4 {
		return 0
	}
	// Simple Laplacian-like: sum of absolute differences with 4-connected neighbors
	var total float64
	var count int
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			idx := y*w + x
			center := gray[idx]
			laplacian := 4*center - gray[idx-1] - gray[idx+1] - gray[idx-w] - gray[idx+w]
			if laplacian < 0 {
				laplacian = -laplacian
			}
			total += laplacian
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// suggestMoods returns mood suggestions based on brightness, contrast, colorfulness.
func suggestMoods(brightness, contrast, colorfulness float64) []string {
	moods := []string{}
	if brightness > 180 {
		moods = append(moods, "bright", "airy", "ethereal")
	} else if brightness < 80 {
		moods = append(moods, "moody", "dramatic", "noir")
	} else {
		moods = append(moods, "balanced", "natural")
	}
	if contrast > 50 {
		moods = append(moods, "dramatic")
	}
	if colorfulness > 100 {
		moods = append(moods, "vibrant")
	} else if colorfulness < 60 {
		moods = append(moods, "subdued")
	}
	// Deduplicate while preserving order
	seen := make(map[string]bool)
	deduped := make([]string, 0, len(moods))
	for _, m := range moods {
		if !seen[m] {
			seen[m] = true
			deduped = append(deduped, m)
		}
	}
	if len(deduped) > 5 {
		deduped = deduped[:5]
	}
	return deduped
}

// suggestStyles returns style suggestions based on orientation, brightness, colorfulness.
func suggestStyles(orientation, brightness, colorfulness string) []string {
	styles := []string{}
	if orientation == "landscape" {
		styles = append(styles, "cinematic wide")
	}
	if brightness == "high-key" {
		styles = append(styles, "soft lighting")
	}
	if brightness == "low-key" {
		styles = append(styles, "chiaroscuro")
	}
	if colorfulness == "vibrant" {
		styles = append(styles, "color grading")
	}
	styles = append(styles, "professional photography")
	if len(styles) > 4 {
		styles = styles[:4]
	}
	return styles
}

// joinComma joins string slices with ", ".
func joinComma(s []string) string {
	if len(s) == 0 {
		return ""
	}
	result := s[0]
	for _, v := range s[1:] {
		result += ", " + v
	}
	return result
}

// joinDot joins string slices with ". ".
func joinDot(s []string) string {
	if len(s) == 0 {
		return ""
	}
	result := s[0]
	for _, v := range s[1:] {
		result += ". " + v
	}
	return result
}
