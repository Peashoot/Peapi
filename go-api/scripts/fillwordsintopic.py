#!/usr/bin/env python2
#-*- coding: UTF-8 -*-
import os
from PIL import Image
from PIL import ImageDraw
from PIL import ImageFont
import matplotlib.pyplot as plt
import numpy as np

def sum(a =1 , b=2):
    print a+b

def fillWordsIntoPic(srd_img_file_path, dst_img_file_path=None, fill_words='try it out', fontFile='', scale=2, sample_step=3):
    try:
        #读取图片信息
        old_img = Image.open(srd_img_file_path)
        pix = old_img.load()
        width = old_img.size[0]
        height = old_img.size[1]

        #创建新图片
        canvas = np.ndarray((height*scale, width*scale, 3), np.uint8)
        canvas[:, :, :] = 255
        new_image = Image.fromarray(canvas)
        draw = ImageDraw.Draw(new_image)
        if fontFile == '':
            abspath = os.path.dirname(os.path.abspath(__file__))
            fontFile = os.path.join(abspath, "source/consola.ttf")
        #创建绘制对象
        font = ImageFont.truetype(fontFile, 10, encoding="unic")
        fill_words = fill_words.decode('utf-8')
        char_table = list(fill_words)

        #开始绘制
        pix_count = 0
        table_len = len(char_table)
        for y in range(height):
           for x in range(width):
               if x % sample_step == 0 and y % sample_step == 0:
                   draw.text((x*scale, y*scale),
                             char_table[pix_count % table_len], pix[x, y], font)
                   pix_count += 1
        # 保存
        if dst_img_file_path is not None:
           new_image.save(dst_img_file_path)
        return 'Success'
    except BaseException as ex:
        return repr(ex)


fillWordsIntoPic(os.path.join(os.path.dirname(
    os.path.abspath(__file__)), "source/timg2.jpg"), "output.jpg", "我的中国♡", os.path.join(os.path.dirname(
        os.path.abspath(__file__)), "source/SimSun.ttf"), 2, 4)
