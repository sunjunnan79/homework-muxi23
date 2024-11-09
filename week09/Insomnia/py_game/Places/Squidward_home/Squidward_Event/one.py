import data
import tool


def A(game):
    data.P = 4


def B(game):
    r = tool.rand(2)
    if r == 1:
        game.update_place("海绵宝宝家")
        game.Text_Printer("很好,海绵宝宝一觉睡到了8点钟,他赶不上早八了(bushi),他火急火燎的冲出了家门")
        data.P = 4
    else:
        game.update_place("海绵宝宝家")
        game.Text_Printer("海绵宝宝躺在床上,当他睁眼时已经是10:00,蟹堡王早就开门。\n当他火急火燎的冲进蟹堡王的时候,只见蟹老板挥舞着大钳子一脸愤怒的看着他,可怜的海绵宝宝...\n")
        game.update_place("蟹堡王")
        data.P = -1


def C(game):
    game.Text_Printer("你坚持不懈，继续去叫章鱼哥。\n（章鱼哥关上了窗户，拉上了窗帘，塞上了耳塞，再也没有回应，你不得不离开,准备前往蟹堡王）\n")
    data.P = 4


def D(game):

    data.P = 2
