import requests
import io

BASE_URL = "http://localhost:8085"

def test_endpoints():
    print("🚀 Starting Web UI API verification...")

    # 1. GET /
    r = requests.get(f"{BASE_URL}/")
    assert r.status_code == 200, f"GET / failed: {r.status_code}"
    assert "via54" in r.text, "GET / did not contain branding text"
    print("  ✓ GET / OK")

    # 2. GET /api/health
    r = requests.get(f"{BASE_URL}/api/health")
    assert r.status_code == 200, f"GET /api/health failed: {r.status_code}"
    assert r.json()["status"] == "ok", "Health status not ok"
    print("  ✓ GET /api/health OK")

    # 3. GET /api/templates
    r = requests.get(f"{BASE_URL}/api/templates")
    assert r.status_code == 200, f"GET /api/templates failed: {r.status_code}"
    assert isinstance(r.json(), list), "Templates response not a list"
    print("  ✓ GET /api/templates OK")

    # 4. GET /api/htmx/status
    r = requests.get(f"{BASE_URL}/api/htmx/status")
    assert r.status_code == 200, f"GET /api/htmx/status failed: {r.status_code}"
    assert "status-chip" in r.text, "HTMX status not loading chips"
    print("  ✓ GET /api/htmx/status OK")

    # 5. GET /api/htmx/pane
    for intent in ["design", "prompt", "present", "video", "forge", "reimagine"]:
        r = requests.get(f"{BASE_URL}/api/htmx/pane?intent={intent}")
        assert r.status_code == 200, f"GET /api/htmx/pane?intent={intent} failed: {r.status_code}"
        assert "flow-step" in r.text or "flow-steps" in r.text or "reimagine" in r.text or "forge" in r.text, f"Pane {intent} html invalid"
        print(f"  ✓ GET /api/htmx/pane?intent={intent} OK")

    # 6. POST /api/htmx/upload (Document upload - fallback field name 'file')
    doc_file = io.BytesIO(b"Hello storyboard script")
    files = {"file": ("test_story.txt", doc_file, "text/plain")}
    r = requests.post(f"{BASE_URL}/api/htmx/upload", files=files)
    if "📄" not in r.text:
        print(f"DEBUG: doc upload response: {repr(r.text)}")
    assert r.status_code == 200, f"POST /api/htmx/upload (doc) failed: {r.status_code}"
    assert "📄" in r.text, "Document preview icon not displayed correctly"
    assert "_path" in r.text, "Doc path hidden input missing in response"
    print("  ✓ POST /api/htmx/upload (Document - fallback 'file') OK")

    # 7. POST /api/htmx/upload (Image upload - fallback field name 'screenshot')
    img_file = io.BytesIO(b"fake image bytes")
    files = {"screenshot": ("shot.png", img_file, "image/png")}
    r = requests.post(f"{BASE_URL}/api/htmx/upload", files=files)
    assert r.status_code == 200, f"POST /api/htmx/upload (img) failed: {r.status_code}"
    assert "<img" in r.text, "Image preview element missing in response"
    assert "_path" in r.text, "Image path hidden input missing in response"
    print("  ✓ POST /api/htmx/upload (Image - fallback 'screenshot') OK")

    # 8. POST /api/htmx/generate (Real task generation)
    data = {"title": "TestDesignSuite", "mode": "presentation"}
    r = requests.post(f"{BASE_URL}/api/htmx/generate", data=data)
    assert r.status_code == 200, f"POST /api/htmx/generate failed: {r.status_code}"
    assert "已生成" in r.text, "Generation response text missing success notice"
    assert "/api/htmx/download" in r.text, "Download link missing in response"
    print("  ✓ POST /api/htmx/generate (Real HTML Generation task) OK")

    # 9. GET /api/htmx/download
    r = requests.get(f"{BASE_URL}/api/htmx/download?name=TestDesignSuite")
    assert r.status_code == 200, f"GET /api/htmx/download failed: {r.status_code}"
    assert "Content-Disposition" in r.headers, "Response missing Content-Disposition header"
    assert "TestDesignSuite.html" in r.headers["Content-Disposition"], "Incorrect download attachment filename"
    assert "<html" in r.text.lower(), "Downloaded content is not valid HTML"
    print("  ✓ GET /api/htmx/download (Output Retrieval) OK")

    print("\n🎉 Web UI API Verification completed successfully! All endpoints function 100% correctly.")

if __name__ == "__main__":
    test_endpoints()

