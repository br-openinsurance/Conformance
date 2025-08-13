import argparse
import os
import json
import re
from .utils import load_allow_list, remove_patch

ALLOWLIST = load_allow_list()

def validate_json(json_obj, api_version_list):
    validation_rules = [
        (len(json_obj) == 4, 'json contains more keys than expected'),
        ('api' in json_obj and remove_patch(json_obj['api']) in api_version_list, 'api field has an invalid value'),
        ('sd' in json_obj and re.match(r'^\d+$', json_obj['sd']), 'sd field has an invalid value'),
        ('docusign_id' in json_obj and re.match(r'^[A-Za-z\d]{8}-[A-Za-z\d]{4}-[A-Za-z\d]{4}-[A-Za-z\d]{4}-[A-Za-z\d]{12}$', json_obj['docusign_id']), 'docusign_id field has an invalid value'),
        ('test_plan_uri' in json_obj and (is_valid_test_plan_uri(json_obj['test_plan_uri'])), 'test_plan_uri field has an invalid value')
    ]

    for rule, message in validation_rules:
        if not rule:
            return False, message

    return True, 'approved'


def is_valid_test_plan_uri(test_plan_uri):
    if isinstance(test_plan_uri, str):
        return bool(re.match(r'^https://web\.conformance\.directory\.opinbrasil\.com\.br/plan-detail\.html\?.*$', test_plan_uri))
    elif isinstance(test_plan_uri, list):
        return all(re.match(r'^https://web\.conformance\.directory\.opinbrasil\.com\.br/plan-detail\.html\?.*$', uri) for uri in test_plan_uri)
    else:
        return False


def parse_args():
    parser = argparse.ArgumentParser(
        description="Checks if the JSON files are correct"
    )
    parser.add_argument(
        "apis",
        nargs="+",
        help="Every API to be checked by the script and their respective versions separated by a '_'. Examples: business_1.3, personal_1.0, consents_2.2."
    )
    return parser.parse_args()


def check_json_files(apis):
    api_list = [api.split('_') for api in apis]
    api_version_list = [api + '_v' + version for api, version in api_list]
    directories = [f"./submissions/functional/{api}/{version}.0" for api, version in api_list]
    
    wrong_files = []

    for directory in directories:
        for filename in os.listdir(directory):
            if filename in ALLOWLIST:
                continue
            if filename.endswith(".json"):
                file_path = os.path.join(directory, filename)
                with open(file_path) as f:
                    try:
                        json_obj = json.load(f)
                        approved, message = validate_json(json_obj, api_version_list)
                        if not approved:
                            wrong_files.append((filename, message))
                    except json.JSONDecodeError as e:
                        wrong_files.append((filename, f'JSON decode error: {str(e)}'))

    return wrong_files


def main():
    args = parse_args()
    wrong_files = check_json_files(args.apis)

    if wrong_files:
        print('The following files contain invalid JSON objects:')
        for filename, message in wrong_files:
            print(f'{filename} ({message})')
        return 1

    return 0


if __name__ == '__main__':
    raise SystemExit(main())
