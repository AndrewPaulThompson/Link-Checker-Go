## Checks the status of \<a> tag links from a url

Install:
```
go install
```
Example:
```
$ link-checker "www.google.com"
[200] - http://www.google.co.uk/imghp?hl=en&tab=wi
[200] - http://maps.google.co.uk/maps?hl=en&tab=wl
[200] - https://play.google.com/?hl=en&tab=w8
[200] - http://www.youtube.com/?gl=GB&tab=w1
[200] - http://news.google.co.uk/nwshp?hl=en&tab=wn
[200] - https://mail.google.com/mail/?tab=wm
[200] - https://drive.google.com/?tab=wo
[200] - https://www.google.co.uk/intl/en/about/products?tab=wh
[200] - http://www.google.co.uk/history/optout?hl=en
[200] - https://accounts.google.com/ServiceLogin?hl=en&passive=true&continue=http://www.google.com/
[200] - http://www.google.com/setprefdomain?prefdom=GB&prev=http://www.google.co.uk/&sig=K_iNPNsUnT2rfYAvHYhA66gC_it3E%3D
```