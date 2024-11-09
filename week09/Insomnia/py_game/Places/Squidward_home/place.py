import Places.Places
import tool
from Places.Squidward_home import Squidward_Event


# 定义地点事件函数
def Start(game):
    Places.Places.CheckSquidward_home  = 1
    # 游戏背景
    background = '''海绵宝宝一大早起来，兴致勃勃地准备叫醒他的好朋友章鱼哥。\n
    他敲了敲章鱼哥的岩石门，喊道：“章鱼哥，快起床！今天是蟹堡王十周年，我们要庆祝呢！”\n
    然而，章鱼哥并不想这么早起来。他在被子里翻了个身，嘟囔着：“海绵宝宝，你知道现在几点吗！？”。\n
    ”4点40分章鱼哥！”海绵宝宝脸上洋溢着微笑。”9:30蟹堡王才开门,还有3个小时!”\n'''

    # 事件添加
    event = {1: '', }

    op = '''请选择以下选项:\n
    A.好吧,显然章鱼哥并不想这么早起床,默默离开,独自前往蟹堡王。\n
    B.章鱼哥说的对，没有必要这么早去蟹堡王，回去睡个回笼觉。\n
    C.你坚持不懈，继续去叫章鱼哥。（章鱼哥关上了窗户，拉上了窗帘，塞上了耳塞，再也没有回应，你不得不离开）\n'''

    ef = {1: {'A': Squidward_Event.one.A, 'B': Squidward_Event.one.B, 'C': Squidward_Event.one.C }}

    if Places.Places.CheckPatrick_home == 0:
        op += 'D.去派大星家（去过了就不再显示）\n'
        ef[1]['D'] = Squidward_Event.one.D

    # 选项添加
    options = {1: op, }

    # 产生的对应影响
    effect = ef

    # 当前的地名
    place = "章鱼哥家"

    # 终极无敌输出器
    tool.Printer(game=game, background=background, event=event, effect=effect, place=place,
                 options=options)
