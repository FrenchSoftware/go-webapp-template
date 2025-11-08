package components

import (
	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
	"github.com/hyperstitieux/template/views/components/icons"
)

// ThemeSwitcher returns a button that toggles between light and dark mode
// Note: Requires theme initialization script in the HTML <head> - see ThemeInitScript()
func ThemeSwitcher() html.Node {
	return html.Button(
		html.Attr("type", "button"),
		html.Attr("aria-label", "Toggle dark mode"),
		html.Attr("onclick", "document.dispatchEvent(new CustomEvent('basecoat:theme'))"),
		attr.Class("btn-icon-outline size-9 flex items-center justify-center"),
		// Both icons are included, CSS classes control which is visible based on theme
		icons.Sun(),
		icons.Moon(),
	)
}

// ThemeInitScript returns the JavaScript code needed to initialize theme switching
// This should be included in the <head> section of your HTML document
func ThemeInitScript() html.Node {
	return html.Raw(`<script>
  (() => {
    try {
      const stored = localStorage.getItem('themeMode');
      if (stored ? stored === 'dark'
                  : matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.classList.add('dark');
      }
    } catch (_) {}

    const apply = dark => {
      document.documentElement.classList.toggle('dark', dark);
      try { localStorage.setItem('themeMode', dark ? 'dark' : 'light'); } catch (_) {}
    };

    document.addEventListener('basecoat:theme', (event) => {
      const mode = event.detail?.mode;
      apply(mode === 'dark' ? true
            : mode === 'light' ? false
            : !document.documentElement.classList.contains('dark'));
    });
  })();
</script>`)
}
