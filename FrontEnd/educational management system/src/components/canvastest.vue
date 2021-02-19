<template>
  <div class="grey lighten-3">
    <v-row>
      <v-col class="col-3 ">
        <v-list shaped elevation="5">
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="title text-center">
                <div class="d-inline-block">选择学生</div>
                <v-progress-circular :value="100*checked/items.length" :width="10" color="teal"
                                     class="float-right d-block" :rotate="-90"></v-progress-circular>
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-virtual-scroll
            :items="items"
            :item-height="100"
            :height="(heights+72)*1.05"
          >
            <template v-slot:default="{item}">
              <v-card
                elevation="2"
                class="ma-4"
                @click="score(item.id)"
                :color="item.isScored==1?'green':'white'"
              >
                <v-card-title class="headline">
                  {{ item.name }}
                </v-card-title>

                <v-card-subtitle>分数：{{ item.score }}</v-card-subtitle>
              </v-card>
            </template>
          </v-virtual-scroll>
        </v-list>
      </v-col>


      <v-col class="col-9">
        <v-card elevation="5">
          <v-container class="mx-auto">
            <v-item-group class="my-9">
              <v-content class="pa-0">
                <v-row class="justify-center">
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on,attrs }">
                      <v-btn dark fab x-large class="mr-3" :color="UIColor"
                             @click.native="showColor = true;AdjustPenSize=false;clearProps();AdjustColor=true"
                             v-bind="attrs" v-on="on">
                        <v-icon>mdi-select-color</v-icon>
                      </v-btn>
                    </template>
                    <span>更换颜色</span>
                  </v-tooltip>
                  <v-snackbar v-model="AdjustColor" timeout="-1" multi-line="true" class="mr-0">
                    <v-row justify="center" class="font-weight-bold my-1">
                      <div style="font-size: 18px">
                        更改颜色
                      </div>
                    </v-row>
                    <v-row justify="center" class="mr-0">
                      <v-btn class="ma-2 pb-16" fab x-large color="red"
                             @click="selectColor('#F44336');sheet=!sheet;UIColor='red';"></v-btn>
                      <v-btn class="ma-2 " fab x-large color="green"
                             @click="selectColor('#4CAF50');sheet=!sheet;UIColor='green';"></v-btn>
                      <v-btn class="ma-2 " fab x-large color="blue"
                             @click="selectColor('#2196F3');sheet=!sheet;UIColor='blue';"></v-btn>
                      <v-btn class="ma-2 " fab x-large color="orange"
                             @click="selectColor('#FF9800');sheet=!sheet;UIColor='orange';"></v-btn>
                      <v-btn class="ma-2 " fab x-large color="pink"
                             @click="selectColor('#E91E63');sheet=!sheet;UIColor='pink';"></v-btn>
                      <v-btn class="ma-2 " fab x-large color="purple"
                             @click="selectColor('#9C27B0');sheet=!sheet;UIColor='purple';"></v-btn>
                      <v-btn class="ma-2 " fab x-large color="black"
                             @click="selectColor('#000000');sheet=!sheet;UIColor='black';"></v-btn>
                    </v-row>
                  </v-snackbar>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="usePen" v-bind="attrs" v-on="on">
                        <v-btn fab dark :color="UIColor" :small="PenIsSmall" :medium="PenIsMedium" :large="PenIsBig"
                               elevation="0"
                               @click="pen;AdjustPenSize=true" @click.native="pen;AdjustPenSize=true" v-bind="attrs"
                               v-on="on">
                          <v-icon dark>mdi-brush</v-icon>
                        </v-btn>
                      </v-btn>
                    </template>
                    <span>使用画笔/修改画笔粗细</span>
                  </v-tooltip>
                  <v-snackbar v-model="AdjustPenSize" timeout="5000" multi-line="true" class="mr-0">
                    <v-row justify="center" class="font-weight-bold my-1">
                      <div style="font-size: 18px">
                        更改画笔大小
                      </div>
                    </v-row>
                    <v-row justify="center" class="mr-0">
                      <v-btn class="mx-2 my-3" fab :color="UIColor"
                             @click="changePenSize(4);PenIsSmall=true;PenIsMedium=false;PenIsBig=false">
                        <v-icon small dark>mdi-brush</v-icon>
                      </v-btn>
                      <v-btn class="mx-2 my-3" fab :color="UIColor"
                             @click="changePenSize(6);PenIsSmall=false;PenIsMedium=true;PenIsBig=false">
                        <v-icon medium dark>mdi-brush</v-icon>
                      </v-btn>
                      <v-btn class="mx-2 my-3" fab :color="UIColor"
                             @click="changePenSize(8);PenIsSmall=false;PenIsMedium=false;PenIsBig=true">
                        <v-icon large dark>mdi-brush</v-icon>
                      </v-btn>
                    </v-row>
                  </v-snackbar>

                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="trueMark();clearProps()"
                             v-bind="attrs"
                             v-on="on">
                        <v-icon dark>mdi-check-bold</v-icon>
                      </v-btn>
                    </template>
                    <span>插入正确符号</span>
                  </v-tooltip>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="errorMark();clearProps()"
                             v-bind="attrs"
                             v-on="on">
                        <v-icon dark>mdi-close-thick</v-icon>
                      </v-btn>
                    </template>
                    <span>插入错误符号</span>
                  </v-tooltip>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="rect();clearProps()" v-bind="attrs"
                             v-on="on">
                        <v-icon dark>mdi-rectangle-outline</v-icon>
                      </v-btn>
                    </template>
                    <span>插入矩形框</span>
                  </v-tooltip>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="arrow();clearProps()"
                             v-bind="attrs"
                             v-on="on">
                        <v-icon dark>mdi-arrow-bottom-left-thick</v-icon>
                      </v-btn>
                    </template>
                    <span>插入箭头</span>
                  </v-tooltip>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click="text();clearProps()" v-bind="attrs"
                             v-on="on">

                        <v-icon dark>mdi-cursor-text</v-icon>
                      </v-btn>
                    </template>
                    <span>插入文字</span>
                  </v-tooltip>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="useEraser" v-bind="attrs"
                             v-on="on">
                        <v-btn fab color="white" :small="EraserIsSmall" :medium="EraserIsMedium" :large="EraserIsBig"
                               elevation="0"
                               @click.native="useEraser">
                          <v-icon dark>mdi-eraser</v-icon>
                        </v-btn>
                      </v-btn>
                    </template>
                    <span>使用橡皮擦/调整橡皮擦大小</span>
                  </v-tooltip>
                  <v-snackbar v-model="AdjustEraserSize" timeout="5000" multi-line="true" class="mr-0" v-bind="attrs"
                              v-on="on">
                    <v-row justify="center" class="font-weight-bold my-1">
                      <div style="font-size: 18px">
                        更改橡皮擦大小
                      </div>
                    </v-row>
                    <v-row justify="center" class="mr-0">
                      <v-btn class="mx-2 my-3" light fab color="white"
                             @click="changeEraserSize(40);EraserIsSmall=true;EraserIsMedium=false;EraserIsBig=false;">
                        <v-icon small>mdi-eraser</v-icon>
                      </v-btn>
                      <v-btn class="mx-2 my-3" light fab color="white"
                             @click="changeEraserSize(60);EraserIsSmall=false;EraserIsMedium=true;EraserIsBig=false;">
                        <v-icon medium>mdi-eraser</v-icon>
                      </v-btn>
                      <v-btn class="mx-2 my-3" light fab color="white"
                             @click="changeEraserSize(80);EraserIsSmall=false;EraserIsMedium=false;EraserIsBig=true;">
                        <v-icon large>mdi-eraser</v-icon>
                      </v-btn>
                    </v-row>
                  </v-snackbar>


                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="undoCanvas();clearProps()"
                             v-bind="attrs" v-on="on">
                        <v-icon dark>mdi-undo</v-icon>
                      </v-btn>
                    </template>
                    <span>撤销</span>
                  </v-tooltip>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="resetCanvas();clearProps()"
                             v-bind="attrs" v-on="on">
                        <v-icon dark>mdi-restart</v-icon>
                      </v-btn>
                    </template>
                    <span>清空笔迹</span>
                  </v-tooltip>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn fab x-large class="mr-3" color="primary" @click.native="saveCanvas();clearProps()"
                             v-bind="attrs" v-on="on">
                        <v-icon dark>mdi-upload</v-icon>
                      </v-btn>
                    </template>
                    <span>上传成绩和批改</span>
                  </v-tooltip>


                </v-row>
              </v-content>
            </v-item-group>
          </v-container>
          <v-container class="mx-auto">
            <canvas id="canvas" :width="widths" :height="heights" class="mx-auto d-block justify-space-between">
              您的浏览器不兼容，请升级或更换浏览器！
            </canvas>
          </v-container>
        </v-card>
      </v-col>
    </v-row>
    <input type="text" value="" class="hide float-right" id="txt" placeholder="请输入文字">
    <div id="res"></div>
    <div class="eraser"></div>
  </div>

</template>

<script>
export default {
  name: "canvastest",
  data() {
    return {
      canvas: {},
      Txt: {},
      ctx: {},
      myCanvas_rect: {},
      imgs: {},
      imgData: [],
      widths: '',
      heights: '',
      color: '#F44336',
      sheet: false,
      UIColor: 'red',
      PenIsSmall: true,
      PenIsMedium: false,
      PenIsBig: false,
      AdjustPenSize: false,
      PenSize: 4,
      AdjustEraserSize: false,
      EraserSize: 20,
      EraserIsSmall: true,
      EraserIsMedium: false,
      EraserIsBig: false,
      AdjustColor: false,
      benched: 0,
      items: [
        {id: 17171020101, name: "zcy", score: 0, isScored: 0},
        {id: 17171020102, name: "zcy1", score: 0, isScored: 0},
        {id: 17171020103, name: "zcy2", score: 0, isScored: 0},
        {id: 17171020104, name: "zcy3", score: 0, isScored: 0},
        {id: 17171020105, name: "zcy4", score: 0, isScored: 0},
        {id: 17171020106, name: "zcy5", score: 0, isScored: 0},
        {id: 17171020107, name: "zcy6", score: 0, isScored: 0},
        {id: 17171020108, name: "zcy7", score: 0, isScored: 0},
        {id: 17171020109, name: "zcy8", score: 0, isScored: 0},
        {id: 17171020110, name: "zcy9", score: 0, isScored: 0},
      ],
      checked: 0,
      windowWidth: 0,
      windowHeight: 0,
    }
  },
  mounted() {
    this.canvas = document.getElementById('canvas');
    this.Txt = document.getElementById('txt');
    this.ctx = canvas.getContext('2d');
    this.imgs = new Image();
    this.imgs.crossOrigin = "Anonymous"
    this.imgs.src = "http://127.0.0.1:9000/image?homework=0000001&student=17171020102";
    this.windowWidth = window.innerWidth - 256
    this.windowHeight = window.innerHeight - 72
    this.imgs.onload = () => {
      if (this.imgs.width > this.windowWidth || this.imgs.height > this.windowHeight) {
        if ((this.imgs.width / this.windowWidth) > (this.imgs.height / this.windowHeight)) {
          this.widths = this.windowWidth * 0.85
          this.heights = (this.imgs.height / this.imgs.width) * this.windowWidth * 0.85
        } else {
          this.heights = this.windowHeight * 0.85
          this.widths = (this.imgs.width / this.imgs.height) * this.windowHeight * 0.85
        }
      } else {
        this.widths = this.imgs.width; //获取图片的宽
        this.heights = this.imgs.height;//获取图片的高
      }

      this.drawImg();
    }
    this.pen();
  },
  methods: {
    changeEraserSize(num) {
      document.querySelector('.eraser').style.display = 'none'
      this.EraserSize = num
    },
    changePenSize(num) {
      this.PenSize = num;
    },
    selectColor(color) {
      this.color = color;
      this.ctx.strokeStyle = color;
      this.showColor = true;
    },
    eraser() {
      var eraser = document.querySelector('.eraser');
      var canvas = this.canvas;
      var ctx = this.ctx;
      this.clearEvent();
      //鼠标按下去,获取鼠标在canvas内的相对位置
      canvas.onmousedown = (ent) => {
        eraser.style.width = this.EraserSize + 'px';
        eraser.style.height = this.EraserSize + 'px';
        eraser.style.position = "absolute"
        eraser.style.display = 'block';

        var scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
        var scrollLeft = document.documentElement.scrollLeft || document.body.scrollLeft;
        canvas.onmousemove = (e) => {
          var oEvent = e || event;
          eraser.style.top = oEvent.clientY + scrollTop  - this.EraserSize / 2 + 'px'
          eraser.style.left = oEvent.clientX + scrollLeft- 56 - this.EraserSize / 2 + 'px'
          ctx.clearRect(e.offsetX, e.offsetY, this.EraserSize, this.EraserSize);
        }
      }
      canvas.onmouseup = function () {
        canvas.onmousemove = null;
        eraser.style.display = 'none'
      }
    },
    drawImg() {
      this.canvas.style.backgroundImage = 'url(' + this.imgs.src + ')'
      this.canvas.style.backgroundRepeat = 'round'
      //this.ctx.drawImage(this.imgs,0,0,this.widths,this.heights);
    },
    rect() {
      var this_ = this;
      var imgData = this.imgData;
      var Txt = this.Txt;
      var ctx = this.ctx;
      var canvas = this.canvas;
      ctx.strokeStyle = this.color;
      this.clearEvent();
      //当前绘制线段的起点
      var startPoint = {};
      //是否进行绘制的开关
      var isDraw = false;
      //保存鼠标拖动之前canvas的数据
      var contextDate = null;
      //绘制的图形，默认是线段
      var drawType = "line";

      function drawHorizontalLine(y) {
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(canvas.width, y);
        ctx.stroke();
      }

      //绘制竖直线，竖直线只有x坐标会变，纵坐标一直是从0到canvas.height
      function drawVerticalLine(x) {
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x, canvas.height);
        ctx.stroke();
      }

      function getLocation(e) {
        var bbox = document.getElementById("canvas").getBoundingClientRect();
        return {
          x: e.offsetX,
          // y: e.offsetY - bbox.top * (canvas.height / bbox.height)
          // x:e.offsetX(),
          y: e.offsetY,
        };
      }

      //绘制矩形方法
      function drawRect(startPoint, endPoint) {
        ctx.save();
        ctx.beginPath();
        ctx.strokeRect(startPoint.x, startPoint.y,
          endPoint.x - startPoint.x, endPoint.y - startPoint.y);
        ctx.stroke();
        ctx.closePath();
      }

      function saveImageData() {
        contextDate = ctx.getImageData(0, 0, canvas.width, canvas.height);
      }

      function restoreImageData() {
        ctx.putImageData(contextDate, 0, 0);
      }

      canvas.onmousedown = function (e) {
        //保存上一次的绘图记录
        imgData.push(ctx.getImageData(0, 0, canvas.width, canvas.height))
        ctx.strokeStyle = this_.color;
        startPoint = getLocation(e);
        isDraw = true;
        saveImageData();
      };
      canvas.onmousemove = function (e) {
        if (isDraw) {
          restoreImageData();
          drawRect(startPoint, getLocation(e));
        }
      };
      canvas.onmouseup = function (e) {
        isDraw = false;
      };
    },
    trueMark() {
      var this_ = this;
      var imgData = this.imgData;
      var Txt = this.Txt;
      var ctx = this.ctx;
      Txt.style.display = "none";
      var canvas = this.canvas;
      ctx.closePath();
      this.clearEvent();
      canvas.onmousedown = function (ev) {
        imgData.push(ctx.getImageData(0, 0, canvas.width, canvas.height))
        ctx.beginPath();
        ctx.strokeStyle = this_.color;
        ctx.lineWidth = 6;
        ctx.lineJoin = "round";

        canvas.onmouseup = function (ev) {
          //鼠标按下的位置
          var oldX = ev.offsetX;
          var oldY = ev.offsetY;

          ctx.moveTo(oldX, oldY);
          ctx.lineTo(ev.offsetX + 20, ev.offsetY + 20);
          ctx.lineTo(ev.offsetX + 60, ev.offsetY - 10);
          ctx.stroke();
        }
        canvas.ondragstart = 'return false';

      };
    },
    errorMark() {
      var this_ = this;
      var imgData = this.imgData;
      var Txt = this.Txt;
      var ctx = this.ctx;
      var canvas = this.canvas;
      Txt.style.display = "none";
      ctx.closePath();
      this.clearEvent();
      canvas.onmousedown = function (ev) {
        imgData.push(ctx.getImageData(0, 0, canvas.width, canvas.height))
        ctx.beginPath();
        ctx.strokeStyle = this_.color;
        ctx.lineWidth = 6;
        ctx.lineJoin = "round";
        //鼠标按下的位置
        canvas.onmouseup = function (ev) {
          var oldX = ev.offsetX;
          var oldY = ev.offsetY;
          ctx.moveTo(oldX, oldY);
          ctx.lineTo(ev.offsetX + 30, ev.offsetY + 30);
          ctx.moveTo(ev.offsetX + 30, ev.offsetY);
          ctx.lineTo(ev.offsetX, ev.offsetY + 30);
          ctx.stroke();
        };
        ctx.closePath();
      };
    },
    pen() {
      var imgData = this.imgData;
      var ctx = this.ctx;
      var canvas = this.canvas;
      this.Txt.style.display = "none";
      var last;
      this.clearEvent();
      // 鼠标按下
      canvas.onmousedown = () => {
        imgData.push(ctx.getImageData(0, 0, canvas.width, canvas.height))
        ctx.lineWidth = this.PenSize;
        ctx.strokeStyle = this.color;
        // 在鼠标按下后触发鼠标移动事件
        canvas.onmousemove = (e) => {
          if (last != null) {
            ctx.beginPath();
            ctx.moveTo(last[0], last[1]);
            ctx.lineTo(e.offsetX, e.offsetY);
            ctx.stroke();
          }
          // 第一次触发这个函数，只做一件事，把当前 鼠标的 x , y 的位置记录下来
          // 做下一次 线段的 起始点。
          last = [e.offsetX, e.offsetY];
        };
      }

      // 鼠标抬起取消鼠标移动的事件
      canvas.onmouseup = () => {
        canvas.onmousemove = null;
        last = null;
      }

      // 鼠标移动函数
      function move(e) {

        if (last != null) {
          ctx.beginPath();
          ctx.lineWidth = pensize;
          ctx.strokeStyle = this_.color;
          ctx.moveTo(last[0], last[1]);
          ctx.lineTo(e.offsetX, e.offsetY);
          ctx.stroke();
        }
        // 第一次触发这个函数，只做一件事，把当前 鼠标的 x , y 的位置记录下来
        // 做下一次 线段的 起始点。
        last = [e.offsetX, e.offsetY];

      }
    },
    arrow() {
      var this_ = this;
      var imgData = this.imgData;
      var Txt = this.Txt;
      var ctx = this.ctx;
      var canvas = this.canvas;
      Txt.style.display = "none";
      var beginPoint = {},
        stopPoint = {},
        polygonVertex = [],
        CONST = {
          edgeLen: 100,
          angle: 50
        };
      this.clearEvent();
      var Plot = {

        angle: "",
        //短距离画箭头的时候会出现箭头头部过大，修改：
        dynArrowSize: function () {
          var x = stopPoint.x - beginPoint.x,
            y = stopPoint.y - beginPoint.y,
            length = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2));
          if (length < 250) {
            CONST.edgeLen = CONST.edgeLen / 2;
            CONST.angle = CONST.angle / 2;
          } else if (length < 500) {
            CONST.edgeLen = CONST.edgeLen * length / 500;
            CONST.angle = CONST.angle * length / 500;
          }
        },
        //getRadian 返回以起点与X轴之间的夹角角度值
        getRadian: function (beginPoint, stopPoint) {
          Plot.angle = Math.atan2(stopPoint.y - beginPoint.y, stopPoint.x - beginPoint.x) / Math.PI * 180;

          paraDef(50, 25);
          Plot.dynArrowSize();
        },
        ///获得箭头底边两个点
        arrowCoord: function (beginPoint, stopPoint) {
          polygonVertex[0] = beginPoint.x;
          polygonVertex[1] = beginPoint.y;
          polygonVertex[6] = stopPoint.x;
          polygonVertex[7] = stopPoint.y;
          Plot.getRadian(beginPoint, stopPoint);
          polygonVertex[8] = stopPoint.x - CONST.edgeLen * Math.cos(Math.PI / 180 * (Plot.angle + CONST.angle));
          polygonVertex[9] = stopPoint.y - CONST.edgeLen * Math.sin(Math.PI / 180 * (Plot.angle + CONST.angle));
          polygonVertex[4] = stopPoint.x - CONST.edgeLen * Math.cos(Math.PI / 180 * (Plot.angle - CONST.angle));
          polygonVertex[5] = stopPoint.y - CONST.edgeLen * Math.sin(Math.PI / 180 * (Plot.angle - CONST.angle));
        },

        //获取另两个底边侧面点
        sideCoord: function () {
          var midpoint = {};
          midpoint.x = (polygonVertex[4] + polygonVertex[8]) / 2;
          midpoint.y = (polygonVertex[5] + polygonVertex[9]) / 2;
          polygonVertex[2] = (polygonVertex[4] + midpoint.x) / 2;
          polygonVertex[3] = (polygonVertex[5] + midpoint.y) / 2;
          polygonVertex[10] = (polygonVertex[8] + midpoint.x) / 2;
          polygonVertex[11] = (polygonVertex[9] + midpoint.y) / 2;
        },

        //画箭头
        drawArrow: function () {
          var canvas = document.getElementById('canvas');
          var ctx = canvas.getContext('2d');
          ctx.beginPath();
          ctx.moveTo(polygonVertex[0], polygonVertex[1]);
          ctx.lineTo(polygonVertex[2], polygonVertex[3]);
          ctx.lineTo(polygonVertex[4], polygonVertex[5]);
          ctx.lineTo(polygonVertex[6], polygonVertex[7]);
          ctx.lineTo(polygonVertex[8], polygonVertex[9]);
          ctx.lineTo(polygonVertex[10], polygonVertex[11]);
          // ctx.lineTo(polygonVertex[0], polygonVertex[1]);
          ctx.closePath();
          ctx.fillStyle = this_.color;
          ctx.fill();
        }
      };
      canvas.onmousemove = null;
      //记录起点beginPoint
      canvas.onmousedown = function (e) {
        imgData.push(ctx.getImageData(0, 0, canvas.width, canvas.height))
        beginPoint.x = e.offsetX;
        beginPoint.y = e.offsetY;
      };

      //记录终点stopPoint，绘图
      canvas.onmouseup = function (e) {
        stopPoint.x = e.offsetX;
        stopPoint.y = e.offsetY;
        Plot.arrowCoord(beginPoint, stopPoint);
        Plot.sideCoord();
        Plot.drawArrow();
      };

      //自定义参数
      function paraDef(edgeLen, angle) {
        CONST.edgeLen = edgeLen;
        CONST.angle = angle;
      }

    },
    text() {
      this.clearEvent();
      var this_ = this;
      var imgData = this.imgData;
      var Txt = null;
      var ctx = this.ctx;
      var canvas = this.canvas;
      var width = this.widths + 200
      canvas.onmousedown = (ev) => {
        if (Txt != null) {
          return;
        }
        oldX = ev.offsetX;
        oldY = ev.offsetY;

        Txt = document.createElement("textArea");
        Txt.style.display = "inline";
        ctx.font = "18px Microsoft Yahei";
        ctx.fontWeight = "800"
        Txt.style.position = 'absolute'
        Txt.style.top = ev.offsetY + 144 + 'px'
        Txt.style.left = ev.offsetX + width / 4 + 'px'

        Txt.style.background = 'rgb(0,0,0,0)'
        Txt.style.border = '3px solid ' + this_.color
        Txt.style.color = this_.color
        Txt.style.width = canvas.width - 20 + "px"

        document.body.append(Txt)
        setTimeout(() => {
          Txt.focus()
        }, "50");
        setTimeout(() => {
          Txt.onblur = fillTxt;
        }, 50);
        canvas.onmousedown = null;

      };
      var oldX = null;
      var oldY = null;

      function fillTxt(ev) {
        imgData.push(ctx.getImageData(0, 0, canvas.width, canvas.height))
        var v = Txt.value;
        if (v !== '') {
          ctx.fillStyle = this_.color
          ctx.moveTo(oldX, oldY);
          let row = []
          let temp = ""
          for (var a = 0; a < v.length; a++) {
            if (ctx.measureText(temp).width < canvas.width - oldX - 20) {
            } else {
              row.push(temp);
              temp = "";
            }
            temp += v[a];
          }
          row.push(temp);
          for (var b = 0; b < row.length; b++) {
            ctx.fillText(row[b], oldX, oldY + (b + 1) * 14);
          }
        }
        Txt.remove();
        Txt = null;
        canvas.onmousedown = (ev) => {
          if (Txt != null) {
            return;
          }
          oldX = ev.offsetX;
          oldY = ev.offsetY;

          Txt = document.createElement("textArea");
          Txt.style.display = "block";
          ctx.font = "16px Microsoft Yahei";
          Txt.style.position = 'absolute'
          Txt.style.top = ev.offsetY + 144 + 'px'
          Txt.style.right = ev.offsetX - ((this.widths + 256) / 4) + 'px'

          Txt.style.background = 'rgb(0,0,0,0)'
          Txt.style.border = '3px solid ' + this_.color
          Txt.style.color = this_.color
          Txt.style.width = canvas.width - oldX - 20 + "px"

          document.body.append(Txt)
          setTimeout(() => {
            Txt.focus()
          }, "50");
          setTimeout(() => {
            Txt.onblur = fillTxt;
          }, 50);
          canvas.onmousedown = null;

        };
      }
    },
    resetCanvas() {
      this.clearEvent();
      var imgData = this.imgData;
      var Txt = this.Txt;
      var ctx = this.ctx;
      var canvas = this.canvas;
      Txt.style.display = "none";
      this.imgData = [];
      ctx.clearRect(0, 0, canvas.width, canvas.height);
    },
    undoCanvas() {
      this.clearEvent();
      var canvas = document.getElementById('canvas');
      var ctx = canvas.getContext('2d');
      var imgData = this.imgData;
      if (imgData.length > 0) {
        let data = imgData.pop()
        // ctx.clearRect(0, 0, canvas.width,canvas.height);
        ctx.putImageData(data, 0, 0, 0, 0, canvas.width, canvas.height);
      }
    },
    saveCanvas() {
      this.clearEvent();
      let tempCanvas = document.createElement("canvas");
      tempCanvas.width = this.widths
      tempCanvas.height = this.heights
      let tempCtx = tempCanvas.getContext('2d');
      if (this.imgData.length === 0) {

        return;
      }
      var image = new Image();
      image = this.canvas.toDataURL('Image/png');
      let imageObj = document.createElement('img');
      imageObj.src = image;
      setTimeout(() => {
        tempCtx.drawImage(this.imgs, 0, 0, this.widths, this.heights);
        tempCtx.drawImage(imageObj, 0, 0, this.widths, this.heights);
        var tempImage = new Image();
        tempImage = tempCanvas.toDataURL('Image/png')
        let arr = tempImage.split(","),
          mime = arr[0].match(/:(.*?);/)[1],
          bstr = atob(arr[1]),
          n = bstr.length,
          u8arr = new Uint8Array(n);
        while (n--) {
          u8arr[n] = bstr.charCodeAt(n);
        }
        let imgFile = new Blob([u8arr], {
          type: mime
        });
        const formData = new FormData();
        formData.append("image", imgFile);
        this.$axios({
          method: "post",
          url: "http://127.0.0.1:9000/image",
          data: formData,
          headers: {
            "Content-Type": "multipart/form-data"
          }
        });
      }, 100);
    },
    clearEvent() {
      var canvas = this.canvas;
      canvas.onmousedown = null;
      canvas.onmousemove = null;
      canvas.onmouseup = null;
    },
    usePen() {
      this.pen();
      this.AdjustPenSize = true;
      this.AdjustEraserSize = false;
      document.querySelector('.eraser').style.display = 'none';
      this.AdjustColor = false;
    },
    useEraser() {
      this.eraser();
      this.AdjustEraserSize = true;
      this.AdjustPenSize = false;
      document.querySelector('.eraser').style.display = 'none'
    },
    clearProps() {
      this.AdjustPenSize = false;
      this.AdjustEraserSize = false;
      document.querySelector('.eraser').style.display = 'none'
      this.AdjustColor = false;
    },
    score(id) {
      this.items.some((item, i) => {
        if (item.id === id) {
          if (this.items[i].isScored === 0) {
            this.items[i].score = 100
            this.items[i].isScored = 1
            this.checked++;
          } else {
            this.items[i].score -= 10
          }
        }
      })
    }
  }
}
</script>

<style scoped>
canvas {
  border: 1px solid black;
}

.hide {
  display: none;
}

#txt {
  position: absolute;
  top: 6%;
  left: 1%;
  width: 575px;
  height: 30px;
  border: 1px solid red;
}

.eraser {
  width: 20px;
  height: 20px;
  border-radius: 100%;
  background: #ccc;
  position: absolute;
  border: 1px solid #666;
  top: 0;
  left: 0;
  display: none;
}

</style>
