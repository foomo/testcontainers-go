import { defineConfig } from 'vitepress'

export default defineConfig({
	base: '/testcontainers-go/',
	title: 'testcontainers-go',
	description: 'Testcontainers integrations for Go',
	themeConfig: {
		logo: '/logo.png',
		outline: [2, 4],
		sidebar: [
			{
				text: 'Overview',
				items: [
					{ text: 'Introduction', link: '/' },
				],
			},
			{
				text: 'Modules',
				items: [
					{ text: 'Temporal', link: '/modules/temporal' },
				],
			},
			{
				text: 'Contributing',
				collapsed: true,
				items: [
					{
						text: "Guideline",
						link: '/CONTRIBUTING.md',
					},
					{
						text: "Code of conduct",
						link: '/CODE_OF_CONDUCT.md',
					},
					{
						text: "Security guidelines",
						link: '/SECURITY.md',
					},
				],
			},
		],
		editLink: {
			pattern: 'https://github.com/foomo/testcontainers-go/edit/main/docs/:path',
			text: 'Suggest changes to this page',
		},
		search: {
			provider: 'local',
		},
		footer: {
			message: 'Made with ♥ <a href="https://www.foomo.org">foomo</a> by <a href="https://www.bestbytes.com">bestbytes</a>',
		},
		socialLinks: [
			{
				icon: 'github',
				link: 'https://github.com/foomo/testcontainers-go',
			},
		],
	},
	head: [
		['meta', { name: 'theme-color', content: '#ffffff' }],
		['link', { rel: 'icon', href: '/logo.png' }],
		['meta', { name: 'author', content: 'foomo by bestbytes' }],
		['meta', { property: 'og:title', content: 'foomo/testcontainers-go' }],
		[
			'meta',
			{
				property: 'og:image',
				content: 'https://github.com/foomo/testcontainers-go/blob/main/docs/public/banner.png?raw=true',
			},
		],
		[
			'meta',
			{
				property: 'og:description',
				content: 'Testcontainers integrations for Go with sane defaults and ergonomic APIs',
			},
		],
		['meta', { name: 'twitter:card', content: 'summary_large_image' }],
		[
			'meta',
			{
				name: 'twitter:image',
				content: 'https://github.com/foomo/testcontainers-go/blob/main/docs/public/banner.png?raw=true',
			},
		],
		[
			'meta',
			{
				name: 'viewport',
				content: 'width=device-width, initial-scale=1.0, viewport-fit=cover',
			},
		],
	],
	markdown: {
		theme: {
			dark: 'github-dark',
			light: 'github-light',
		}
	},
	sitemap: {
		hostname: 'https://foomo.github.io/testcontainers-go',
	},
	ignoreDeadLinks: true,
})
