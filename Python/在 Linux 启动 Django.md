# 在 Linux 启动 Django

## User:
请问如何在linux中启动django？

## ChatGPT:
在 Linux 中启动 Django 项目，你可以按照以下的步骤来操作：

1. 安装 Python：Django 是 Python 语言写的，因此你需要先安装 Python。可以使用 apt-get (Ubuntu/Debian) 或 yum (CentOS) 命令来安装。

```
sudo apt-get install python3
```

2. 安装 Pip：Pip 是 Python 的包管理器，用于安装和管理 Python 包。可以使用下面的命令安装：

```
sudo apt-get install python3-pip
```

3. 安装 Django：现在你可以使用 pip 安装 Django 了。

```
pip3 install Django
```

4. 创建 Django 项目：现在你可以创建一个新的 Django 项目。替换 `myproject` 为你想要的项目名。

```
django-admin startproject myproject
```

5. 启动 Django 项目：导航到你的项目目录，然后运行 Django 开发服务器。替换 `myproject` 为你的项目名。

```
cd myproject
python3 manage.py runserver
```


默认情况下，这将启动一个运行在 localhost 的 8000 端口的开发服务器。你可以在浏览器中访问 http://localhost:8000  查看你的网站。

注意：Django 开发服务器仅用于开发环境，不应在生产环境中使用。在生产环境中，你应该配置一个更加健壮的 WSGI 服务器，例如 Gunicorn 或 uWSGI。

Thread generated with <a href="https://chat.openai.com">ChatGPT</a> and downloaded with <a href="https://botrush.io">Botrush</a>