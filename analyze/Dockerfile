FROM ubuntu

# get up pip, vim, etc.
RUN apt-get -y update --fix-missing
RUN apt-get install -y python-pip python-dev libev4 libev-dev gcc libxslt-dev libxml2-dev libffi-dev vim curl 
RUN pip install --upgrade pip

# get numpy, scipy, scikit-learn and flask
RUN apt-get install -y python-numpy python-scipy python-matplotlib
RUN pip install pandas

# add our project
ADD analyze.py /analyze.py
