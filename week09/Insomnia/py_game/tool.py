import random
import shutil
import threading

from rich.console import Console
from rich.panel import Panel
from rich.live import Live
from rich.text import Text
from rich.table import Table
import time

import data

console = Console()


# 终极无敌文本输出器
class TextGame:
    def __init__(self):
        # 创建Console对象，用于输出到终端
        self.console = Console()

        # 创建四个面板对象，用于显示对话、历史记录、选项和地点
        self.dialogue_panel = Panel("Dialogue", title="[b]Dialogue[/b]")
        self.history_panel = Panel("History", title="[b]History[/b]")
        self.options_panel = Panel("Options", title="[b]Options[/b]")
        self.place_panel = Panel("Place", title="[b]Place[/b]")  # 修改此处

        # 计算并设置面板布局，使每个面板占据终端大小的四分之一
        self.calculate_layout()
        self.layout = Table.grid(expand=True)
        self.layout.add_row(self.dialogue_panel, self.place_panel)  # 修改此处
        self.layout.add_row(self.history_panel, self.options_panel)

        # 创建Live对象，用于实时在终端中显示输出
        self.live = Live(self.layout, console=self.console)
        self.is_running = False  # 用于控制主循环是否运行

    def start(self):
        # 开始主循环
        self.is_running = True
        threading.Thread(target=self.run).start()

    def stop(self):
        # 停止主循环
        self.is_running = False

    def calculate_layout(self):
        # 获取终端的宽度和高度
        terminal_width, terminal_height = shutil.get_terminal_size()

        # 计算每个面板应该占据的宽度和高度
        panel_width = terminal_width // 2
        panel_height = terminal_height // 2

        # 设置每个面板的宽度和高度
        self.dialogue_panel.width = panel_width
        self.dialogue_panel.height = panel_height
        self.history_panel.width = panel_width
        self.history_panel.height = panel_height
        self.options_panel.width = panel_width
        self.options_panel.height = panel_height
        self.place_panel.width = panel_width
        self.place_panel.height = panel_height

    def update_dialogue(self, dialogue):
        # 更新对话面板的内容
        self.dialogue_panel.renderable = Text(dialogue, style="bold magenta")
        self.refresh()

    def update_history(self, history):
        # 更新历史记录面板的内容
        self.history_panel.renderable = Text(history, style="italic yellow")
        self.refresh()

    def update_options(self, options):
        # 更新选项面板的内容
        self.options_panel.renderable = Text(options, style="bold green")
        data.options = options
        self.refresh()

    def update_place(self, place):  # 修改方法名
        # 更新视频面板的内容
        self.place_panel.renderable = Text(place, style="underline blue")
        self.refresh()

    def refresh(self):
        # 每次循环前重新计算面板布局
        self.calculate_layout()
        time.sleep(0.001)  # 调整刷新频率，单位为秒，根据需要进行调整
        self.live.refresh()

    def run(self):
        # 在Live对象中开始主循环，实时刷新输出
        with self.live:
            while True:
                # 每次循环前重新计算面板布局
                self.calculate_layout()
                time.sleep(0.001)  # 调整刷新频率，单位为秒，根据需要进行调整
                self.live.refresh()

    # 用于统一的在当前剧情上面的输出
    def Text_Printer(self, text=''):
        history = ''
        # 初始化当前对话
        dialogues = ''
        for char in text:
            dialogues += char
            time.sleep(0.1)  # 调整刷新频率，单位为秒，根据需要进行调整
            self.update_dialogue(dialogues)
            if char == '\n' or len(dialogues)==46:
                if len(dialogues)==46:
                    dialogues+='\n'
                # 输出到历史文本区
                history += dialogues
                self.update_history(history)
                # 清空当前的对话文本
                dialogues = ''
    # 选择判断器,让用户输入选择,不为对应范围就让他重新输入
    def Choice(self,choice_num):
        # 启动选择,顺便自动把他转化为大写
        choice = input().upper()

        # 如果不按照格式输出就让他重新输出
        while choice > chr(ord('A') + choice_num - 1) or choice < 'A':
            self.options_panel.renderable = Text(data.options+'\n输入格式错误,请输入A~'+chr(ord('A') + choice_num - 1), style="bold green")
            choice = input().upper()
        return choice


# 移动花费时间函数,返回了移动完的时间,和本次移动花费的实际时间
def Pass_Time(pass_time):
    # use_time =速率 * 事件正常花费的时间
    use_time = data.speed * pass_time
    # 更改你当前的时间
    data.time += use_time
    # 返回当前的时间和你花费的实际时间
    return data.time, use_time


# 随机数控制器
def rand(n):
    return random.randint(1, n) % (n + 1)





def Printer(game, background, event, effect, place, options, next_place=-100,end = ''):
    # 更新地点
    game.update_place(place)

    # 产生一个随机数,用来进行随机事件
    r = rand(len(event))

    # 游戏背景
    game.Text_Printer(background + event[r])

    # 输出对应选项
    game.update_options(options[r])

    # 调用相关函数,后面是用来获取选择的函数
    effect[r][game.Choice(len(effect[r]))](game)

    if end != '':
        game.Text_Printer(end)

    if next_place != -100:
        # 更改游戏的下一个地点
        data.P = next_place
