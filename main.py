import sys
import subprocess
import json
from PyQt5.QtWidgets import QMessageBox, QFileDialog, QMainWindow, QApplication
from PyQt5.Qt import QIntValidator
from UI_main import Ui_MainWindow

import os

EXECUTABLE_PATH = None


class potplayer:
    def __init__(self):
        self.playlist = []
        self.executable, _ = self.find_executable()

    def output_to_file(self, stream_list: list, path: str):
        with open(path, "w") as f:
            f.write(self.list_2_dpl(stream_list))

    def run_playlist(self, playlist: list, tmpPath="./output.dpl"):
        """Open a playlist/a video/a audio/a image with PotPlayer.
        Python will pause while potplayer is playing.
        """
        res = self.list_2_dpl(playlist)

        f = open(tmpPath, "w")
        f.write(res)
        f.close()

        cmd = '"%s" "%s"' % (self.executable, tmpPath)
        subprocess.Popen(cmd)

    def set_executable(self, path):
        self.executable = path

    def find_executable(self):
        """Automatically find PotPlayer executable file path and process name.
        It depends on your system.
        """
        x86 = r"C:\Program Files (x86)\DAUM\PotPlayer\PotPlayerMini.exe"
        x64 = r"C:\Program Files\DAUM\PotPlayer\PotPlayerMini64.exe"
        if os.path.exists(x86):
            return x86, os.path.basename(x86)
        elif os.path.exists(x64):
            return x64, os.path.basename(x64)
        else:
            try:
                process_name = os.path.basename(EXECUTABLE_PATH)
                if process_name not in ["PotPlayerMini.exe", "PotPlayerMini64.exe"]:
                    raise ValueError("Cannot find potplayer executable! "
                                     "Please edit '%s' to add the valid path.")
                return EXECUTABLE_PATH, process_name
            except:
                raise ValueError("Cannot find potplayer executable! "
                                 "Please edit '%s' to add the valid path.")

    def list_2_dpl(self, stream_list: list):
        res = "DAUMPLAYLIST\n"
        lines = list()
        for index in range(len(stream_list)):
            lines.append("{0}*file*{1}".format(index + 1, stream_list[index]))
            lines.append("{0}*played*0".format(index + 1))
        res += "\n".join(lines)
        return res


class window(QMainWindow):
    def __init__(self):
        super().__init__()
        self.ui = Ui_MainWindow()
        self.ui.setupUi(self)

        self.select_interface = str(self.ui.comboBox_selectInterface.currentIndex() + 1)
        self.room_realID = ""
        self.quality = []
        self.stream = []
        self.potplayer = potplayer()

        self.validator = QIntValidator(self)
        self.ui.lineEdit_truthId.setValidator(self.validator)
        self.ui.lineEdit.setValidator(self.validator)
        self.ui.lineEdit_truthId.setReadOnly(True)
        self.ui.lineEdit_potPlayerPath.setText(self.potplayer.executable)

        self.ui.comboBox_selectInterface.currentIndexChanged.connect(self.change_interface)
        self.ui.pushButton_getQuality.clicked.connect(self.btn_get_realId_quailty)
        self.ui.pushButton_getStream.clicked.connect(self.btn_get_stream)
        self.ui.pushButton.clicked.connect(self.btn_output_file)
        self.ui.pushButton_2.clicked.connect(self.btn_open_potplayer)
        self.ui.pushButton_setPotplayer.clicked.connect(self.btn_setPotPlayer)

    def change_interface(self, index):
        self.select_interface = str(index + 1)

    def btn_get_realId_quailty(self, checked):
        roomId = self.ui.lineEdit.text().strip()
        if roomId == "":
            QMessageBox.warning(self, "错误", "请输入直播间号")
        else:
            res = get_realId_quality(roomId, self.select_interface)

            if res["type"] != 0:
                QMessageBox.warning(self, "错误", res["data"])
            else:
                self.ui.lineEdit_truthId.setText(res["realId"])
                self.room_realID = res["realId"]
                self.quality = []
                self.ui.comboBox_selectQuality.clear()
                self.ui.listWidget.clear()
                for qua in res["Quality"]["quality"]:
                    self.quality.append(qua)
                    self.ui.comboBox_selectQuality.addItem(qua["desc"])

    def btn_get_stream(self, checked):
        real_id = self.ui.lineEdit_truthId.text()
        if real_id == "":
            QMessageBox.warning(self, "错误", "请先读取直播间")
        else:
            self.ui.listWidget.clear()
            res = get_stream(real_id,
                             str(self.quality[self.ui.comboBox_selectQuality.currentIndex()]["qn"]),
                             self.select_interface)
            self.stream = res["urls"]
            for url in res["urls"]:
                self.ui.listWidget.addItem(url)

    def btn_output_file(self, checked):
        if len(self.stream) == 0:
            QMessageBox.warning(self, "错误", "请先获取直播流")
        else:
            fileName = QFileDialog.getSaveFileName(self, "Save File",
                                                   "./output.dpl",
                                                   "File (*.dpl)")
            if fileName[0] == '':
                return
            else:
                self.potplayer.output_to_file(self.stream, fileName[0])

    def btn_open_potplayer(self, checked):
        self.potplayer.run_playlist(self.stream)

    def btn_setPotPlayer(self, checked):
        name, _ = QFileDialog.getOpenFileName(self, "选择文件", "./", "可执行文件 (*.exe)")
        if name == '':
            return
        self.ui.lineEdit_potPlayerPath.setText(name)


def get_realId_quality(live_id, apiType):
    # res = {
    #     "type": 0,   ##succ
    #     "Quality": [{"qn": 250}, {"desc": "高清"}, ],
    #     "realId": 123156
    # }
    ret = subprocess.run(["src.exe", "-id", live_id, "-type", "0", "-apiType", apiType],
                         stdout=subprocess.PIPE)
    res = json.loads(ret.stdout)
    return res


def get_stream(real_id: str, qn: str, apiType: str):
    # res = {
    #     urls:["","", ""]
    # }
    ret = subprocess.run(["src.exe", "-id", real_id, "-type", "1", "-quality", qn, "-apiType", apiType],
                         stdout=subprocess.PIPE)
    return json.loads(ret.stdout)


if __name__ == '__main__':
    app = QApplication(sys.argv)
    window = window()
    window.setWindowTitle("bilibili获取直播间流")
    window.show()
    sys.exit(app.exec_())
