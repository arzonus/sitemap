# SiteMap
SiteMap is application for scanning sites and building site map. 
It creates map of urls on site and subsites.. 

# Install
`go get -u github.com/arzonus/sitemap/cmd/sitemap`

# Run
After installation you can use the app

```
$ sitemap --depth=1 https://facebook.com

2019/05/18 15:55:25 **********************
2019/05/18 15:55:25 * url      https://facebook.com
2019/05/18 15:55:25 * file     sitemap.out
2019/05/18 15:55:25 * depth    0
2019/05/18 15:55:25 * timeout  30s
2019/05/18 15:55:25 * parallel 12
2019/05/18 15:55:25 **********************
2019/05/18 15:55:25 
https://facebook.com
├── https://www.facebook.com/
├── https://www.facebook.com/recover/initiate?lwv=110&ars=royal_blue_bar
├── https://en-gb.facebook.com/
├── https://uk-ua.facebook.com/
├── https://fi-fi.facebook.com/
├── https://zh-cn.facebook.com/
├── https://de-de.facebook.com/
├── https://ar-ar.facebook.com/
├── https://tr-tr.facebook.com/
├── https://fr-fr.facebook.com/
├── https://es-la.facebook.com/
├── https://pt-br.facebook.com/
├── https://messenger.com/
├── https://l.facebook.com/l.php?u=https%3A%2F%2Finstagram.com%2F&h=AT3W8rMmsZFLXYn9pZRHi3F6zEPtyxEC3T2f-j-eogsV3BlTx1a2XkGYWvEUXRq-5Z1AG8aYFdTHlwWQ3MfxpTD8WexgNyZ_tRMzbL_Qz5xUTm51nupHdBNc4Xwxyrji0B5-XP2GcsMff2z_
├── https://developers.facebook.com/?ref=pf
├── https://www.facebook.com/help/568137493302217
├── https://www.facebook.com/help/2687943754764396
└── https://www.facebook.com/help/www/1573156092981768/
2019/05/18 15:55:25 time: 463.340153ms
2019/05/18 15:55:25 result saved to file sitemap.out
```