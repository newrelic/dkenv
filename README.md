DKENV
-----
dkenv is a tool that downloads Docker versions for you, keeps track of your versions, and lets you switch between docker versions on the fly. You can also specify the apiversion on the command line and dkenv will select the correct docker version.

Designed to bypass the dreaded:

    2014/08/26 14:21:03 Error response from daemon: client and server don't have same version (client : 1.13, server: 1.12)


Use:

    dkenv -version version

    dkenv -apiVersion apiVersion

Version example: "1.6.0"

API Version example: "1.18"

dkenv stores the docker files in ~/.dkenv and creates a symlink in /usr/local/bin

Contributions
-------------

Contributions are more than welcome. Bug reports with specific reproduction
steps are great. If you have a code contribution you'd like to make, open a
pull request with suggested code.

Pull requests should:

 * Clearly state their intent in the title
 * Have a description that explains the need for the changes
 * Include tests!
 * Not break the public API
 * Add yourself to the CONTRIBUTORS file. I might forget.

If you are simply looking to contribute to the project, taking on one of the
items in the "Future Additions" section above would be a great place to start.
Ping us to let us know you're working on it by opening a GitHub Issue on the
project.

By contributing to this project you agree that you are granting New Relic a
non-exclusive, non-revokable, no-cost license to use the code, algorithms,
patents, and ideas in that code in our products if we so choose. You also agree
the code is provided as-is and you provide no warranties as to its fitness or
correctness for any purpose

Copyright (c) 2015 New Relic, Inc. All rights reserved.