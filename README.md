# heat-solver

# Please install this tools 
pip install pandas matplotlib numpy

# How to Run : 
## To run the project with front 
go run cmd/server/main.go

## To run the project without front 
go run cmd/head/main.go --method=FTCS --dx=0.1 --dt=0.0005 --tmax=1.0 --out=ftcs.csv
go run cmd/head/main.go --method=BTCS --dx=0.1 --dt=0.0005 --tmax=1.0 --out=btcs.csv
go run cmd/head/main.go --method=CN --dx=0.1 --dt=0.0005 --tmax=1.0 --out=cn.csv

### To visualize 
python plot_results.py compare ftcs.csv FTCS btcs.csv BTCS cn.csv CN

