import tool
import Places.Places
from Places.Patrick_home import Patrick_home_Event


# 定义地点事件函数
def Start(game):
    Places.Places.CheckPatrick_home = 1
    # 游戏背景
    background = '''你敲了敲派大星的家门,没有任何动静,很遗憾派大星可能还在睡觉,毕竟时间还很早。\n'''
    # 事件添加
    event = {1: '', }

    op = '请选择以下选项:\nA:再敲一下\nB:默默离开\n'
    ef = {1: {'A': Patrick_home_Event.one.A, 'B':Patrick_home_Event.one.B}}

    if Places.Places.CheckSquidward_home == 0:
        op += 'C:去章鱼哥家\n'
        ef[1]['C'] = Patrick_home_Event.one.C

    # 选项添加
    options = {1: op, }

    # 产生的对应影响
    effect = ef

    # 当前的地名
    place = "派大星家"

    # 终极无敌输出器
    tool.Printer(game=game, background=background, event=event, effect=effect, place=place,
                 options=options)
