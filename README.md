# filtermanuel
Filters a Missing Manuel file by what is available in a faxbot, grouped by KoL area

This script is intended to filter a file by lines contained in another file
This is as trivial as it sounds; we need to maintain insertion order
More practically (in terms of Kingdom of Loathing):
* Take the output from a monster manuel status utility
* Find the matching monsters in a faxbot list
* Output the list of matches (in manuel order) into an output file
* Preserve area headings from manuel in the output