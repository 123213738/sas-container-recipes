#!/bin/bash -e
#
# Copyright 2018 SAS Institute Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Enable debugging if SAS_DEBUG is set
[[ -z ${SAS_DEBUG+x} ]] && export SAS_DEBUG=0
if [ ${SAS_DEBUG} -gt 0 ]; then
    set -x
fi

if [[ -z ${RUN_USER+x} ]]; then
    if [[ -z ${CASENV_ADMIN_USER+x} ]]; then
        export RUN_USER="cas"
    else
        export RUN_USER="${CASENV_ADMIN_USER}"
    fi
fi

echo "[INFO] : RUN_USER set to ${RUN_USER}"

if [[ -z ${JUPYTER_TOKEN+x} ]]; then
    echo "[WARNING] : JUPYTER_TOKEN is not provided so blank will be used."
fi

RUN_USER_HOME=$(getent passwd ${RUN_USER} | cut -d: -f6)
if [ ! -d "${RUN_USER_HOME}" ]; then
    echo
    echo "[INFO] : Creating ${RUN_USER}'s home directory of ${RUN_USER_HOME}"
    mkdir --verbose --parents ${RUN_USER_HOME}
    chown --verbose ${RUN_USER} ${RUN_USER_HOME}
    chmod --verbose 0700 ${RUN_USER_HOME}
    echo
else
    echo
    echo "[INFO] : ${RUN_USER}'s home directory of ${RUN_USER_HOME} already exists"
    echo
fi

###############################################################################
# Run Jupyter notebook
###############################################################################
[[ -z ${ENABLE_TERMINAL+x} ]]      && ENABLE_TERMINAL=True
[[ -z ${ENABLE_NATIVE_KERNEL+x} ]] && ENABLE_NATIVE_KERNEL=True

export JPY_COOKIE_SECRET=`openssl rand -hex 32`
export AUTHINFO="${RUN_USER_HOME}/authinfo.txt"
export SSLCALISTLOC="${SASHOME}/SASSecurityCertificateFramework/cacerts/trustedcerts.pem"
export CAS_CLIENT_SSL_CA_LIST="/data/casconfig/sascas.pem"

cp /usr/local/lib/python3.6/site-packages/saspy/sascfg.py /usr/local/lib/python3.6/site-packages/saspy/sascfg_personal.py
sed -i -e "s#/opt/sasinside/SASHome/SASFoundation/9.4/bin/sas_u8#/opt/sas/spre/home/SASFoundation/sas#g" \
    /usr/local/lib/python3.6/site-packages/saspy/sascfg_personal.py

_jupyterpid="/var/run/jupyter.pid"
touch ${_jupyterpid}

runuser --shell "/bin/sh" --login ${RUN_USER} \
    --command "mkdir -p --verbose ~/jupyter"

# In the following, the echo was added after switching to Python 3.6
# http://forums.fast.ai/t/jupyter-notebook-fails-to-start/8370/5
echo
echo "[INFO] : Create jupyter config file"
echo
runuser --shell "/bin/sh" --login ${RUN_USER} \
    --command "if [ ! -e \"${RUN_USER_HOME}/.jupyter/jupyter_notebook_config.py\" ]; then jupyter notebook --generate-config; echo \"c.NotebookApp.allow_remote_access = True\" >> ~/.jupyter/jupyter_notebook_config.py; fi"

# This is bad, but pause for a moment to make sure the generate call has completed
# There have been some cases where the image just stops after the above call.
# It does not happen every time.

sleep 3

echo
echo "[INFO] : Starting jupyter..."
echo

runuser --shell "/bin/sh" --login ${RUN_USER} \
    --command "JPY_COOKIE_SECRET=${JPY_COOKIE_SECRET} \
    AUTHINFO=${AUTHINFO} \
    SSLCALISTLOC=${SSLCALISTLOC} \
    CAS_CLIENT_SSL_CA_LIST=${CAS_CLIENT_SSL_CA_LIST} \
    jupyter notebook \
    --ip='*' \
    --no-browser \
    --NotebookApp.token='${JUPYTER_TOKEN}' \
    --NotebookApp.terminals_enabled=${ENABLE_TERMINAL} \
    --NotebookApp.base_url=/Jupyter \
    --KernelSpecManager.ensure_native_kernel=${ENABLE_NATIVE_KERNEL} \
    --notebook-dir=~/jupyter &"
pgrep jupyter > ${_jupyterpid}
