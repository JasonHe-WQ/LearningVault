import subprocess

command = "test -e /home/hewenqing/images."
try:
    result = subprocess.check_call(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
except subprocess.CalledProcessError as e:
    print("error", e)
    print("return code", e.returncode)
else:
    print(result == 0)

print("result", "check_all.py")
# print(result.stdout)
# print(result.returncode)
# print(len(result.communicate()))
