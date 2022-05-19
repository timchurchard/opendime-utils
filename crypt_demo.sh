#!/bin/bash
./opendime-utils crypt --address 133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg --message "Hello World" --signature "Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU=" -e -i "Test Message for crypt" -o
#Encrypted message for 133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg
#BEWs+6LFlRxUQQ4X9kuPp03L3C+qttMjliLWVExapsDQdjaFfV/7sPxHbhVDPeZ2upYx99TzK0TufWEupSUAXC7s69dbdqUTiTYZOkfKRahrpJNmTffwbmgIO+lI8qNk/SBXVR/CNu+toq/H+5KqJ6njeqaNZX8=

./opendime-utils crypt -i "BEWs+6LFlRxUQQ4X9kuPp03L3C+qttMjliLWVExapsDQdjaFfV/7sPxHbhVDPeZ2upYx99TzK0TufWEupSUAXC7s69dbdqUTiTYZOkfKRahrpJNmTffwbmgIO+lI8qNk/SBXVR/CNu+toq/H+5KqJ6njeqaNZX8=" -k L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK -d -o
#Decrypted message
#Test Message for crypt

./opendime-utils crypt --address 133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg --message "Hello World" --signature "Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU=" -e --inputfile crypt_demo.sh --outputfile crypt_demo.sh.enc
./opendime-utils crypt -d -k L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK --inputfile crypt_demo.sh.enc --outputfile crypt_demo.sh.dec
