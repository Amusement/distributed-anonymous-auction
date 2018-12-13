import subprocess
import paramiko

from constants import *

result = subprocess.run([(
        "az vmss list-instance-public-ips --resource-group %s --name %s --query [].{ip:ipAddress} -o tsv" % (
    ResourceGroup, VMSSname))], stdout=subprocess.PIPE, shell=True)
publicIPs = result.stdout.decode("utf-8").split("\n")[0:-1]

commands = []

for ip in publicIPs[1:]:
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    ssh.connect(ip, username=VMusername, password=VMpassword)

    ssh_stdin, ssh_stdout, ssh_stderr = ssh.exec_command("sudo go run P2-d3w9a-b3c0b-b3l0b-k0b9/auctioneer_main.go P2-d3w9a-b3c0b-b3l0b-k0b9/auctioneer/config.json")
    commands.append((ssh_stdout, ssh_stderr, ssh))

for command in commands:
    ssh_stdout, ssh_stderr, ssh = command
    ssh_stdout.channel.recv_exit_status()

    lines = ssh_stdout.readlines()
    for line in lines:
        print(line) 

    lines = ssh_stderr.readlines()
    for line in lines:
        print(line)

    ssh.close()

