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

    print('File count in ', 'directory:  ' + dir, file_count - 1, ', directory size is: ', dir_size)


def get_file_marker(root_path):
    dir = root_path + 'occurrences'
    list = os.listdir(dir)
    print(list)


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
    file_marker_input = input('Input file marker to restore: ')
    count_number_and_size(get_root_path(), 'occurrences/' + file_marker_input)

    start_restore = datetime.datetime.now()
    test_restore_file(get_root_path(), file_marker_input)
    finish_restore = datetime.datetime.now()
    work_time_restore = finish_restore - start_restore
    full_work_time = work_time_save + work_time_restore
    time_list.append(full_work_time.total_seconds())
    print('Time for saving file is ', work_time_save)
    print('Time for restoring file is ', work_time_restore)

    delete_from_dir(get_root_path() + '/batches')
    delete_from_dir(get_root_path() + '/occurrences')



def make_plot_work_time(time_list):
    x = ['md5', 'sha1', 'sha256', 'sha512']
    plt.bar(x, time_list, label='Seconds')
    plt.xlabel('Algoritms')
    plt.ylabel('Time')
    plt.title('Time spent on execution')
    plt.legend()
    plt.show()


def make_plot_losses():
    x = ['md5', 'sha1', 'sha256', 'sha512']
    plt.bar(x, diff_list, label='Bytes')
    plt.xlabel('Algorithms')
    plt.ylabel('Number of Losses')
    plt.title('Losses after restore')
    plt.legend()
    plt.show()


def bitwise_compare(file1, file2):
    count = 0
    with open(file1, "rb") as f1, open(file2, "rb") as f2:
        while True:
            byte1 = f1.read(1)
            byte2 = f2.read(1)

            if not byte1 and not byte2:
                print("The files are identical.")
                break
            elif not byte1 or not byte2:
                print("Files vary in size.")
                break
            elif byte1 != byte2:
                print(f"Files vary per byte: {byte1} (in file {file1}) != {byte2} (in file {file2})")
                count = count + 1
    print('Number of different bytes', count)
    return count


if __name__ == "__main__":
    PATH_TO_FILE = input('Input File Path to save: ')
    diff_list = []
    time_list = []
    run_dedup('md5', time_list, diff_list)

    restored_file_path = input('Input restored file: ')
    diff_list.append(bitwise_compare(PATH_TO_FILE, restored_file_path))

    run_dedup('sha1', time_list, diff_list)

    diff_list.append(bitwise_compare(PATH_TO_FILE, restored_file_path))

    run_dedup('sha256', time_list, diff_list)

    diff_list.append(bitwise_compare(PATH_TO_FILE, restored_file_path))

    run_dedup('sha512', time_list, diff_list)

    diff_list.append(bitwise_compare(PATH_TO_FILE, restored_file_path))

    make_plot_work_time(time_list)
    make_plot_losses()
    print(diff_list)
