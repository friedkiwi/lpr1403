# lpr1403

A (semi) drop-in replacement for lpr to receive output through https://1403.bitnet.systems

# Usage

Just use instead of lpr. The following command line flags are available:

 - `-auth-token` The token you've received from virtual1403
 - `-endpoint` The endpoint to submit the jobs to. By default, the bitnet.systems one is being used.

When using it with lp5250d to connect to an AS/400 or IBM i machine, the following template can be used:

    lp5250d -N env.DEVNAME=MACDESK outputcommand='scs2ascii | ./lpr1403 -auth-token "authtokenData=" ' as400host.net

# Build instructions

    go build github.com/cyberdotgent/lpr1403
