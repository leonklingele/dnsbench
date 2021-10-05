# How-to update domains.txt

1. Visit https://www.alexa.com/topsites
2. Run the following JavaScript:

```js
{
	let topsites = [];

	document.querySelectorAll('.site-listing').
		forEach((e) => {
			const a = e.querySelector('a');
			const domain = a.getAttribute('href').replace('/siteinfo/', '');
			topsites.push(domain);
		});

	console.log(topsites.join('\n'));
}
```

