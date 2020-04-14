# Chaos Toolkit for Litmus Chaos

## ChaosToolKit

The Chaos Toolkit aims to be the simplest and easiest way to explore building your own Chaos Engineering Experiments. It also aims to define a vendor and technology independent way of specifying Chaos Engineering experiments by providing an Open API.

Reference: https://chaostoolkit.org/

## Steps to install Chaos Toolkit

Install Python for your system:

1. On MacOS X:
   ```
   $ brew install python3
   ```

1. On Debian/Ubuntu:
   ```
   $ sudo apt-get install python3 python3-venv
   ```

1. On CentOS:
   ```
   $ sudo yum -y install https://centos7.iuscommunity.org/ius-release.rpm
   $ sudo yum -y install python35u
   ```
   > **Note:**, on CentOS, the Python 3.5 binary is named python3.5 rather than python3 as other systems.

1. On Windows:
   ```
   Download the latest binary installer (https://www.python.org/downloads/windows/) from the Python website.
   ```

# Local Development

1. In this directory
    ```
    cd test-tools/chaostoolkit
   ```
1. build python package
    ```
    python setup.py develop
   ```
1. In this directory
    ```
    cd test-tools
    ```
1. build pip module
    ```
    pip install chaostoolkit/
   ```


