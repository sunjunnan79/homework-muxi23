from Places import Places
import data
import Places

def A(game):
    count = 0
    l = 2
    op = "请选择以下选项:\nA.再敲一下\nB.默默离开\n"
    if Places.Places.CheckSquidward_home == 0:
        op += 'C:去章鱼哥家\n'
        l += 1

    options = op
    # 更新选项
    game.Text_Printer("周围十分寂静,并没有发生任何事\n")
    game.update_options(options)
    # 进行选择
    choice = game.Choice(l)
    #假如地点没发生转变的话
    while (data.P == 2):
        #如果选了B
        if choice == 'B':
            data.P = 3
        #如果选了C
        elif choice == 'C':
            data.P = 1
        #如果选了A
        else:
            game.Text_Printer("周围十分寂静,并没有发生任何事\n")
            game.update_options(options)

            if count == 4:
                game.Text_Printer("派大星大声吼叫:谁这么没素质啊,大早上的在这大吼大叫,滚远一点!!!\n")
                op = "请选择以下选项:\nA.默默离开\n"
                if Places.Places.CheckSquidward_home == 0:
                    op += 'B:去章鱼哥家\n'
                    l += 1
                game.update_options(op)
                choice = game.Choice(l - 2)
                if choice == 'A':
                    #去驾校
                    data.P = 3
                else:
                    #去章鱼哥家
                    data.P = 1
            else:
                count+=1
                choice = game.Choice(l)


def B(game):
    data.P = 4
    return



def C(game):
    data.P = 2
