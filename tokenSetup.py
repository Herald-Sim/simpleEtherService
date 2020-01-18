"""
All codes from 'OpenZepplin-contracts' repository
repository: https://github.com/OpenZeppelin/openzeppelin-contracts
# ERC20
# https://raw.githubusercontent.com/OpenZeppelin/openzeppelin-contracts/master/contracts/token/ERC20/ERC20.sol
# IERC20
# https://raw.githubusercontent.com/OpenZeppelin/openzeppelin-contracts/master/contracts/token/ERC20/IERC20.sol
# SafeMath
# https://raw.githubusercontent.com/OpenZeppelin/openzeppelin-contracts/master/contracts/math/SafeMath.sol
# GSN
# https://raw.githubusercontent.com/OpenZeppelin/openzeppelin-contracts/master/contracts/GSN/Context.sol
"""

import wget
import os

ERC20URL = "https://raw.githubusercontent.com/OpenZeppelin/openzeppelin-contracts/master/contracts/token/ERC20/ERC20.sol"
IERC20URL = "https://raw.githubusercontent.com/OpenZeppelin/openzeppelin-contracts/master/contracts/token/ERC20/IERC20.sol"
SafeMathURL = "https://raw.githubusercontent.com/OpenZeppelin/openzeppelin-contracts/master/contracts/math/SafeMath.sol"
GSN = "https://raw.githubusercontent.com/OpenZeppelin/openzeppelin-contracts/master/contracts/GSN/Context.sol"

# Download latest contract code
print("Start get latest ERC20 contract code\n")
fileERC20 = wget.download(ERC20URL)
fileIERC20 = wget.download(IERC20URL)
fileSafeMath = wget.download(SafeMathURL)
fileGSN= wget.download(GSN)
print("\nGet latest ERC20 complete\n")


# Modify ERC20.sol code (import path)
safeMathPath = 'import "./SafeMath.sol";'
gsnPath = 'import "./Context.sol";'

ERC20sol = open("./ERC20.sol", "r+")
ERC20solBuffer = []

# read all file contents custom buffer(ERC20solBuffer)
while True:
    line = ERC20sol.readline()
    if not line:
        break
    ERC20solBuffer.append(line)

# change import path
ERC20solBuffer[4] = safeMathPath
ERC20solBuffer[2] = gsnPath
ERC20sol.close()

os.system("rm ./ERC20.sol")

# create new ERC20.sol file
newERC20sol = open("./ERC20.sol", "w+")

for line in ERC20solBuffer:
    newERC20sol.write(line)

newERC20sol.close()

print("ERC20.sol modify complete")
