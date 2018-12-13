import subprocess
import threading

import paramiko

from constants import *

result = subprocess.run([(
        "az vmss list-instance-public-ips --resource-group %s --name %s --query [].{ip:ipAddress} -o tsv" % (
    ResourceGroup, VMSSname))], stdout=subprocess.PIPE, shell=True)
publicIPs = result.stdout.decode("utf-8").split("\n")[0:-1]



class SSHThread(threading.Thread):
    def __init__(self, ip):
        super(SSHThread, self).__init__()
        self.ip = ip


    def run(self):
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(self.ip, username=VMusername, password=VMpassword)

        ssh_stdin, ssh_stdout, ssh_stderr = ssh.exec_command(
            "sudo go run P2-d3w9a-b3c0b-b3l0b-k0b9/auctioneer_main.go P2-d3w9a-b3c0b-b3l0b-k0b9/auctioneer/config.json")

        def line_buffered(f):
            line_buf = ""
            while not f.channel.exit_status_ready():
                line_buf += f.read(1)
                if line_buf.endswith('\n'):
                    yield line_buf
                    line_buf = ''

        for l in line_buffered(ssh_stdout):
            print(self.ip+": "+ l)
            
        ssh.close()


threads = []
for ip in publicIPs[1:]:
    ssh_thread = SSHThread(ip)
    threads.append(ssh_thread)
    ssh_thread.start()


for thread in threads:
   thread.join()

