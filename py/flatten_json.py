import json

# coding by v2er @Projection
def unwind(root):
    def dfs(node, inherit=False):
        state = states[-1]
        if not inherit:
            states.append(state := states[-1].copy())
        for k, v in node.items():
            if isinstance(v, list):
                continue
            path.append(k)
            if isinstance(v, dict):
                dfs(v, True)
            else:
                state[".".join(path)] = v
            path.pop()

        is_parent = any(isinstance(v, list) for v in node.values())
        if is_parent:
            for k, v in node.items():
                if isinstance(v, list):
                    for vv in v:
                        path.append(k + "[*]")
                        dfs(vv)
                        path.pop()
        elif not inherit:
            result.append(state.copy())
        if not inherit:
            states.pop()

    path = ["$"]
    states = [{}]
    result = []
    dfs(root, True)
    return result


def main():
    with open("../data.json", "rb") as f:
        data = json.load(f)

    got = unwind(data)
    print("--- GOT:")
    print(json.dumps(got, indent='\t'))

    with open("../want.json", "rb") as f:
        want = json.load(f)

    if got == want:
        print("--- PASS")
    else:
        print("--- FAIL")


if __name__ == "__main__":
    main()
