import data
import tool


def A(game):
    return


def B(game):
    # 具体逻辑自己写,随便你怎么用,但是请一定用game里面的方法
    options = "请选择以下选项:\nA.去章鱼哥家\nB.去派大星家\n"
    # 更新选项
    game.update_options(options)
    # 进行选择
    choice = game.Choice(2)
    if choice == 'A':
        data.P = 1
    else:
        data.P = 2

