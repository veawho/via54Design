package template

type LayoutTemplate struct {
	ID       string     `yaml:"id"`
	Name interface{} `yaml:"name"`
	Version  string     `yaml:"version"`
	Category string     `yaml:"category"`
	Tags     []string   `yaml:"tags"`
	When     LayoutWhen `yaml:"when"`
	Elems    []Element  `yaml:"elements"`
	CSS      string     `yaml:"css"`
}

type LayoutWhen struct {
	ContentHas  []string `yaml:"content_has"`
	SuitableFor []string `yaml:"suitable_for"`
}

type Element struct {
	Role     string    `yaml:"role"`
	Position string    `yaml:"position,omitempty"`
	Behavior string    `yaml:"behavior,omitempty"`
	Tag      string    `yaml:"tag,omitempty"`
	Style    string    `yaml:"style,omitempty"`
	FontSize string    `yaml:"font_size,omitempty"`
	Children []Element `yaml:"children,omitempty"`
}

type PaletteItem struct {
	Role        string `yaml:"role"`
	Hex         string `yaml:"hex"`
	NameZh      string `yaml:"name_zh,omitempty"`
	CulturalNote string `yaml:"cultural_note,omitempty"`
}

type ColorSchemeTemplate struct {
	ID           string            `yaml:"id"`
	Name interface{} `yaml:"name"`
	Version      string            `yaml:"version"`
	Tags         []string          `yaml:"tags"`
	Source       string            `yaml:"source,omitempty"`
	Colors       map[string]string `yaml:"colors,omitempty"`
	Palette      []PaletteItem     `yaml:"palette,omitempty"`
	When         ColorWhen         `yaml:"when"`
	CSSVariables string            `yaml:"css_variables,omitempty"`
	Mood         []string          `yaml:"mood,omitempty"`
	Season       string            `yaml:"season,omitempty"`
}

type ColorWhen struct {
	BrandTone   []string `yaml:"brand_tone"`
	Audience    []string `yaml:"audience"`
	SuitableFor []string `yaml:"suitable_for"`
}

type TypographyTemplate struct {
	ID          string            `yaml:"id"`
	Name        interface{}       `yaml:"name"`
	Version     string            `yaml:"version"`
	Tags        []string          `yaml:"tags"`
	Fonts       map[string]string `yaml:"fonts"`
	Weights     map[string]string `yaml:"weights,omitempty"`
	Sizes       map[string]string `yaml:"sizes"`
	GoogleFonts []string          `yaml:"google_fonts,omitempty"`
	Mood        []string          `yaml:"mood,omitempty"`
	When        TypographyWhen    `yaml:"when"`
}

type TypographyWhen struct {
	BrandTone   []string `yaml:"brand_tone"`
	SuitableFor []string `yaml:"suitable_for"`
}

type TemplateRegistry struct {
	Layouts      []RegistryEntry `yaml:"layouts"`
	ColorSchemes []RegistryEntry `yaml:"color_schemes"`
	Typography   []RegistryEntry `yaml:"typography"`
	Narratology  []RegistryEntry `yaml:"narratology"`
	Combinations []Combination   `yaml:"combinations"`
}

type RegistryEntry struct {
	ID       string   `yaml:"id"`
	Name interface{} `yaml:"name"`
	Version  string   `yaml:"version"`
	Category string   `yaml:"category"`
	Tags     []string `yaml:"tags"`
	File     string   `yaml:"file"`
}

type Combination struct {
	Name interface{} `yaml:"name"`
	Layout  string   `yaml:"layout"`
	Color   string   `yaml:"color"`
	Font    string   `yaml:"font"`
	Suitable []string `yaml:"suitable"`
}

type GenerationResult struct {
	HTML         string
	LayoutID     string
	ColorID      string
	FontID       string
	Title        string
	FontImports  string
	CSSVars      string
	BaseCSS      string
	LayoutCSS    string
	LetteringSVG string // 手写/书法 SVG path (可选)
}
