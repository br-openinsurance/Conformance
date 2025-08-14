import re


ALLOWLIST_FILEPATH = "old_file_names.txt"
_VERSION_TRIPLET = re.compile(r'(?<!\d)(\d+)\.(\d+)\.(\d+)(?!\d)')


def remove_patch(api_version: str) -> str:
    return _VERSION_TRIPLET.sub(r'\1.\2', api_version, count=1)

def load_allow_list():
    with open(ALLOWLIST_FILEPATH, "r") as file:
        return file.read().splitlines()

if __name__ == "__main__":
    print(load_allow_list())