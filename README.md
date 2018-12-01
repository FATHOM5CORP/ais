<p align="left">
  <img src="https://static1.squarespace.com/static/596d24cd4402430bb863ffad/t/5b41e62603ce641f98f2e3cd/1536741696061" width="100" title="HACKtheMACHINE Seattle">
</p>

Installation... `go get github.com/FATHOM5/ais`

# AIS *Beta Release*
In September of 2018 the United States Navy hosted their annual [HACKtheMACHINE](https://hackthemachine.ai) Navy Digital Experience in Seattle, Washington.  The public competition was centered on three competitive tracks with prizes for top finishers.  **Track 2 of the competition was Data Science and the Seven Seas: Collision Avoidance.**   

The U.S. Navy explained to participants that they are the largest international operator of unmanned and autonous systems sailing on and under the world's oceans and that there is a common interest with the Navy and the public to design algorithms that contribute to safa navigation by autonomous machines piloted by artifiical intelligence.  Additionally, the U.S. Navy is a advocate and enforcer of the international right to freedom of the seas.  Extending that freedom to safe navigation by autonous systems is a natural extension of long-standing traditions. To support the development of such AI driven navigational systems the Navy sponsored HACKtheMACHINE Seattle Track 2 in order to begin developing datasets from publicly available maritime shipping data in order to eventually train a machine learning algorithm to make prudent decisions for avoiding a collision at sea.  Read the full challenge description [here](https://github.com/FATHOM5/Seattle_Track_2).

This repository is a Go language package and an open-source release  of the winning solutions to HACKtheMACHINE Seattle, Track2: Data Science and the Seven Seas.  [FATHOM5](https://fathom5.co) is a contract partner of the U.S. in executing HACKtheMACHINE and is proud to have developed this repo as a contribution to the maritime data science community.  **We are growing an inclusive community of data science practitioners in support of maritime issues.  Please use this code, submit issues to improve it and join in.  Our community is organized here and on LinkedIn [here](https://www.linkedin.com/groups/12145028/). Please reach out with questions about the code, maritime data science or just to ask a few questions to learn more.** 

## What's in the package?
*Package `FATHOM5/AIS` contains tools for creating machine learning datasets for autonomous navigation system development based on open data released by the U.S. Government.*

The largest and most comprehensive public data source for maritime domain awareness is the Automatic Identification System (AIS) data collected and released to the public by the U.S. Government on the [marinecadastre.gov](https://marinecadastre.gov/ais/) website.  These comma separated value (csv) data files average more than 25,000,000 records per file, and a single month of data is a set of 20 files totalling over 60Gb of information.  Therefore, the first challenge in building a machine learning dataset from these files is a big-data challenge to find interesting interactions in this large corpus of records.

The package contains tools for abstracting the process opening and reading these large files and additional tools that support algorithm development to identify interesting interactions.  See *Usage* below for additional details.

##Installation
Package `FATHOM5/ais` is a standard Go language library installed in the typical fashion.

    go get github.com/FATHOM5/ais
    




