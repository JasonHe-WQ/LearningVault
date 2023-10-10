from collections import deque

# 现在,你有三个槽位,和对应的法力,当槽位为冰冰冰时,消耗法力100,冰冰火消耗175,冰冰雷消耗200,冰火火75,冰火雷300,冰雷雷150,
# 火火火175,火火雷200,火雷雷60,雷雷雷125.请问给定初始状态,能够最多召唤多少种技能？你可以控制加入冰火雷,但是槽位只有三个,
# 如冰冰冰,加入雷,则为雷冰冰.输入第一个参数arg1为列表,包含三个数字,代表当前状态,0为空,可以随意添加,1,2,3代表冰雷火,
# arg2为法力值,arg3为最大按键次数,加入元素需要按键,释放技能需要按键

def max_skills(arg1, arg2, arg3):
    # Skill costs
    skill_costs = {
        '111': 100,
        '113': 175,
        '112': 200,
        '133': 75,
        '132': 300,
        '122': 150,
        '333': 175,
        '332': 200,
        '322': 60,
        '222': 125,
    }

    # Initialize DP table
    # Key is a tuple (mana, state), value is the maximum number of skills that can be cast
    dp = {(arg2, ''.join(map(str, arg1))): 0}

    # Initialize queue for BFS, each element is a tuple (remaining mana, state, max skills so far)
    queue = deque([(arg2, ''.join(map(str, arg1)), 0)])

    max_skills = 0

    while queue and arg3 > 0:
        next_queue = deque()

        while queue:
            mana, state, skills = queue.popleft()

            # Try casting each skill
            for skill, cost in skill_costs.items():
                if mana >= cost and state.endswith(skill[:-1]):
                    next_mana = mana - cost
                    next_skills = skills + 1
                    next_state = state[1:] + skill[-1]

                    # Update DP table and queue if this is a better way to reach this state
                    if (next_mana, next_state) not in dp or dp[(next_mana, next_state)] < next_skills:
                        dp[(next_mana, next_state)] = next_skills
                        next_queue.append((next_mana, next_state, next_skills))

                    max_skills = max(max_skills, next_skills)

            # Try pressing each button (1=ice, 2=thunder, 3=fire)
            for button in '123':
                next_state = state[1:] + button
                if (mana, next_state) not in dp or dp[(mana, next_state)] < skills:
                    dp[(mana, next_state)] = skills
                    next_queue.append((mana, next_state, skills))

        queue = next_queue
        arg3 -= 1

    return max_skills

# Test the function
arg1 = [0, 0, 0]  # initial state: empty slots
arg2 = 1000  # initial mana
arg3 = 10  # maximum number of button presses
print(max_skills(arg1, arg2, arg3))
