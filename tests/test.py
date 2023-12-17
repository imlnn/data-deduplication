import json
import os
import datetime
import shutil
import matplotlib.pyplot as plt


def get_root_path():
    full_path = os.getcwd()
    root_path = full_path.removesuffix('tests')
    return root_path


def test_save_file(root_path, file_path):
    os.system('cd ' + root_path + '&& ./dedup save ' + file_path + '>null 2>&1')


def test_restore_file(root_path, file_marker):
    os.system('cd ' + root_path + '&& ./dedup restore ' + file_marker + '>null 2>&1')


def count_number_and_size(root_path, file_dir):
    file_count = 0
    dir = root_path + file_dir

    for root, dirs, files in os.walk(dir):
        file_count += len(files)

    dir_size = sum(d.stat().st_size for d in os.scandir(dir) if d.is_file())

    print('File Count in ', 'Directory:  ' + dir, file_count - 1, ', Directory Size is: ', dir_size)


def get_file_marker(root_path):
    dir = root_path + 'occurrences'
    list = os.listdir(dir)
    print(list)


def get_file_size(path_to_restored_file, path_to_saved_file):
    saved_file_stats = os.stat(path_to_saved_file)
    restored_file_stats = os.stat(path_to_restored_file)
    print(f'Saved File Size in Bytes is {saved_file_stats.st_size}')
    print(f'Restored File Size in Bytes is {restored_file_stats.st_size}')
    return saved_file_stats.st_size - restored_file_stats.st_size


def delete_from_dir(directory_path):
    try:
        with os.scandir(directory_path) as entries:
            for entry in entries:
                if entry.is_file():
                    os.unlink(entry.path)
                else:
                    shutil.rmtree(entry.path)
        print("All files and subdirectories deleted successfully.")
    except OSError:
        print("Error occurred while deleting files and subdirectories.")


def change_alg_config(alg):
    with open(get_root_path() + 'config.json', 'r') as f:
        config = json.load(f)

    config['alg'] = alg

    with open(get_root_path() + 'config.json', 'w') as f:
        json.dump(config, f)


def run_dedup(alg, time_list, diff_list):
    change_alg_config(alg)

    start_save = datetime.datetime.now()
    test_save_file(get_root_path(), PATH_TO_FILE)
    finish_save = datetime.datetime.now()
    work_time_save = finish_save - start_save

    count_number_and_size(get_root_path(), 'batches')

    get_file_marker(get_root_path())
    file_marker_input = input('Input File Marker to restore: ')
    count_number_and_size(get_root_path(), 'occurrences/' + file_marker_input)

    start_restore = datetime.datetime.now()
    test_restore_file(get_root_path(), file_marker_input)
    finish_restore = datetime.datetime.now()
    work_time_restore = finish_restore - start_restore
    full_work_time = work_time_save + work_time_restore
    time_list.append(full_work_time.total_seconds())
    print('Time for Saving File is ', work_time_save)
    print('Time for Restoring File is ', work_time_restore)
    difference_file_size = get_file_size(input('Input Restored File: '), PATH_TO_FILE)
    diff_list.append(difference_file_size)

    delete_from_dir(get_root_path() + '/batches')
    delete_from_dir(get_root_path() + '/occurrences')


def make_plot(time_list):
    x = ['md5', 'sha1', 'sha256', 'sha512']
    plt.bar(x, time_list, label='Seconds')
    plt.xlabel('Algoritms')
    plt.ylabel('Time')
    plt.title('Time spent on execution')
    plt.legend()
    plt.show()


if __name__ == "__main__":
    PATH_TO_FILE = input('Input File Path to save: ')
    diff_list = []
    time_list = []
    run_dedup('md5', time_list, diff_list)
    run_dedup('sha1', time_list, diff_list)
    run_dedup('sha256', time_list, diff_list)
    run_dedup('sha512', time_list, diff_list)
    make_plot(time_list)
    print(diff_list)
