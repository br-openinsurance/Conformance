ALLOWLIST_FILEPATH = "old_file_names.txt"

def load_allow_list():
    with open(ALLOWLIST_FILEPATH, "r") as file:
        return file.read().splitlines()

if __name__ == "__main__":
    print(load_allow_list())