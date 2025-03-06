@echo off
setlocal enabledelayedexpansion

:: Function to prompt the user for input
set "PromptUser="
set /p PromptUser="Do you want to make a new commit? (y/n): "
if /i "!PromptUser!"=="y" (
    set "commitMessage="
    set /p commitMessage="Enter the commit message: "
    git commit -am "!commitMessage!"
)

:: Perform interactive rebase
set "numberOfCommits="
set /p numberOfCommits="Enter the number of commits to squash: "
git rebase -i HEAD~!numberOfCommits!

:: Ask if the user wants to force push with lease
set "forcePush="
set /p forcePush="Do you want to force push with lease? (y/n): "
if /i "!forcePush!"=="y" (
    git push --force-with-lease
)

endlocal