import atexit
import Places.Bob_home.place as Bob_home
import data as data
import Places.Patrick_home.place as Patrick_home
import Places.Squidward_home.place as Squidward_home
from tool import TextGame

# 事件添加器,把你写的地点函数添加到对应的位置
place = {0: Bob_home.Start, 1: Squidward_home.Start, 2: Patrick_home.Start}

CheckPatrick_home = 0
CheckSquidward_home = 0


def Start_Game():
    # 创建TextGame对象
    game = TextGame()
    game.start()
    # 让游戏一直进行下去根据游戏的位置触发游戏事件,当P没有被设置为-1,游戏就会根据位置P的变化自动切换到对应为止
    while data.P != -1:
        place[data.P](game)
    # 确保程序退出时停止线程
    atexit.register(game.stop)
