import tool
from Places.Bob_home import Bob_home_Event


# 定义地点事件函数
def Start(game):
    
    # 游戏背景
    background = '''今天是蟹堡王十周年，海绵宝宝一大早就从床上里跃出来，眼睛里闪烁着期待和兴奋。他的方形裤子已经穿好，脸上里还残留着早餐的泡沫。\n
现在，他站在菠萝屋的门口，准备出发前往蟹堡王。阳光透过海水，照在他的脸上，让他感到温暖和愉快。\n他知道，今天的上班之路一定会是一场奇妙的冒险。于是，他迈开步伐，蓝色的身影在海底大道上跳跃着，向着蟹堡王的方向前进。\n
这是一个充满惊喜和欢乐的日子，而海绵宝宝的心里充满了期待。\n他不知道，今天会发生什么奇妙的事情，但他知道，无论发生什么，他都会用他的微笑和善良来面对。\n'''

    # 事件添加
    event = {1: '', }

    # 选项添加
    options = {1: '请选择以下选项:\nA:直接出发前往蟹堡王\nB:先去和邻居们打招呼\n', }

    # 产生的对应影响
    effect = {1: {'A': Bob_home_Event.one.A, 'B': Bob_home_Event.one.B, }, }

    # 当前的地名
    place = "海绵宝宝家"

    # 终极无敌输出器
    tool.Printer(game=game, background=background, event=event, effect=effect, place=place,
                 options=options)
