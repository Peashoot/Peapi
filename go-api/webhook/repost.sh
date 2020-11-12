#!/bin/bash
cd /root/code/TreeHole
hugo
rm -rf public
git update
./deploy.sh
cd /usr/share/nginx
rm -rf blog
cp /root/code/TreeHole/public .
mv public blog
cd /root/code/TreeHole/public
git init
git config --global user.email "1106328900@qq.com" 
git config --global user.name "peashoot"
git remote add origin git@github.com:Peashoot/Peashoot.github.io.git
git add .
git commit -m "deploy"
git push -u origin master