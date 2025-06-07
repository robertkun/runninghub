1. 调用runninghub后台api接口, 实现Api的调用, 已经工作流的调用, 完成自动生图和生视频的功能
2. your-api-key: e5ed5c9ad9234ed185ec07b46022309b

RunningHub 原生 ComfyUI 接口支持说明
https://www.runninghub.cn/proxy/{your-api-key}

工作流:
1930266544381792258



调用刚才生成的图片上传api, 上传本地图片, test.png, 然后再调用任务api, 
生成一个任务, 把上传的图片, 做为参数, 传给任务, 调用的workflow id 为 1930266544381792258, 

然后, 修改节点18的参数, 设置为刚才上传的图片地址

{
      "id": 18,
      "type": "LoadImage",
      "pos": [
        -562.9784545898438,
        227.427001953125
      ],
      "size": [
        576.5151977539062,
        1259.9168701171875
      ],
      "flags": {
        "pinned": true
      },
      "order": 14,
      "mode": 0,
      "inputs": [],
      "outputs": [
        {
          "label": "IMAGE",
          "name": "IMAGE",
          "type": "IMAGE",
          "slot_index": 0,
          "links": [
            43,
            378
          ]
        },
        {
          "label": "MASK",
          "name": "MASK",
          "type": "MASK"
        }
      ],
      "properties": {
        "Node name for S&R": "LoadImage"
      },
      "widgets_values": [
        "e8242f2541d974adccafd0fd774981a2556a5265b659c61699d60a57a56134de.png",
        "image",
        ""
      ],
      "color": "#bec7b6",
      "bgcolor": "#aab3a2"
    },