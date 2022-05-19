#!/bin/bash
./crypt --address 133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg --message "Hello World" --signature "Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU=" -e -i "Test Message for crypt" -o
#Encrypted message for 133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg
#BEWs+6LFlRxUQQ4X9kuPp03L3C+qttMjliLWVExapsDQdjaFfV/7sPxHbhVDPeZ2upYx99TzK0TufWEupSUAXC7s69dbdqUTiTYZOkfKRahrpJNmTffwbmgIO+lI8qNk/SBXVR/CNu+toq/H+5KqJ6njeqaNZX8=

./crypt -i "BEWs+6LFlRxUQQ4X9kuPp03L3C+qttMjliLWVExapsDQdjaFfV/7sPxHbhVDPeZ2upYx99TzK0TufWEupSUAXC7s69dbdqUTiTYZOkfKRahrpJNmTffwbmgIO+lI8qNk/SBXVR/CNu+toq/H+5KqJ6njeqaNZX8=" -k L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK -d -o
#Decrypted message
#Test Message for crypt

./crypt --address 133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg --message "Hello World" --signature "Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU=" -e --inputfile demo.sh --outputfile demo.sh.enc
./crypt -d -k L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK --inputfile demo.sh.enc --outputfile demo.sh.dec
