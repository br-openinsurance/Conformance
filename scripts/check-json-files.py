import argparse
import os
import json
import re


def validate_json(json_obj, api_version_list):
    if len(json_obj) != 4:
        return False, 'json contains more keys than expected'
    if 'api' not in json_obj or json_obj['api'] not in api_version_list:
        return False, 'api field has an invalid value'
    if 'sd' not in json_obj or not re.match(r'^\d+$', json_obj['sd']):
        return False, 'sd field has an invalid value'
    if 'docusign_id' not in json_obj or not re.match(r'^[A-Z\d]{8}-[A-Z\d]{4}-[A-Z\d]{4}-[A-Z\d]{4}-[A-Z\d]{12}$', json_obj['docusign_id']):
        return False, 'docusign_id field has an invalid value'
    if 'test_plan_uri' not in json_obj or not re.match(r'^https://web\.conformance\.directory\.opinbrasil\.com\.br/plan-detail\.html\?.*$', json_obj['test_plan_uri']):
        return False, 'test_plan_uri field has an invalid value'
    return True, 'approved'


def main(argv = None):
    parser = argparse.ArgumentParser(
        description="Checks if the json files are correct"
    )
    parser.add_argument(
        "apis",
        nargs="+",
        help="Every api to be checked by the script and their respective versions separated by a '_'. Examples: business_1.3, personal_1.0, consents_2.2."
    )
    args = parser.parse_args(argv)

    api_list = [api.split('_') for api in args.apis]
    api_version_list = [api + '_v' + version for api, version in api_list]
    directories = [f"./submissions/functional/{api}/{version}.0" for api, version in api_list]
    
    wrong_files = []
    for directory in directories:
        for filename in os.listdir(directory):
            if filename.endswith(".json"):
                file_path = os.path.join(directory, filename)
                with open(file_path) as f:
                    json_obj = json.load(f)
                    approved, message = validate_json(json_obj, api_version_list)
                    if not approved:
                        wrong_files.append((filename, message))

    if len(wrong_files):
        print('The following files contain invalid json objects:')
        for file in wrong_files:
            filename, message = file
            print(f'{filename} ({message})')
        return 1

    return 0


if __name__ == '__main__':
    raise SystemExit(main())
