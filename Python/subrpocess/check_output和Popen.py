import subprocess
import logging

logging.basicConfig(level=logging.INFO)
command = "test -e /home/hewenqing/images.txt && echo 'True' || echo 'False'"
logging.info(command)
check_output = subprocess.check_output(command, shell=True)
# .check_output可以直接读取，读出来是byte类型，需要decode成utf-8
print("check", check_output)
check_output_bool = check_output.strip().decode('utf-8') == 'True'
print("check_output_bool", check_output_bool)
command2 = ["ls /home/hewenqing/images.txt"]
popen_output = subprocess.Popen(command2, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
print('Command: ', command2)
print("popen stdout", popen_output.stdout)
# 如果这一行读取，使用popne_output.stdout.read()，会导致下面的popen_output.communicate()读出为空
print("return code", popen_output.returncode)
print("popen communicate", popen_output.communicate())
# .communicate()方法会阻塞，直到子进程结束，返回一个元组(stdoutdata, stderrdata)
print("popen terminate")
popen_output.terminate()
# .terminate()方法会立即终止子进程
print("return code", popen_output.returncode)
print("popen communicate", popen_output.communicate())
output_bool = popen_output == 'True'
print(output_bool)
