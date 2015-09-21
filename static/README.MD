###Setup

For the UI to work it is necessary to copy the content of [sample nginx config file](/docs/nginx/nginx.default.linux.txt) in your nginx default file, usually
```
/etc/nginx/sites-available/default
```
Change the path accordingly to match your environment and restart nginx service.
```
sudo service nginx restart
```
Once nginx changes are done, you will need to configure grunt. From ```static/``` directory just run
```
npm install
bower install
```
This will download the required UI dependencies. The final step is to create the ```dist/``` files which are minimized version for production.

You can just run ```grunt``` and it will deploy to ```dist/``` directory with the default css theme, which is ```sandstone``` from [bootswatch](http://bootswatch.com/sandstone/)
```
grunt
```
Or you can pass which theme to copy to ```dist/``` directory like this
```
grunt --theme=yeti
```
When you're done with UI, browse to your [http://localhost/](http://localhost/) and the application should be ready.
