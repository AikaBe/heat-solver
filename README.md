# heat-solver

# Сначала решаем задачу всеми методами
go run cmd/head/main.go --method=FTCS --dx=0.1 --dt=0.0005 --tmax=1.0 --out=ftcs.csv
go run cmd/head/main.go --method=BTCS --dx=0.1 --dt=0.0005 --tmax=1.0 --out=btcs.csv
go run cmd/head/main.go --method=CN --dx=0.1 --dt=0.0005 --tmax=1.0 --out=cn.csv

# Теперь сравниваем визуально
python plot_results.py compare ftcs.csv FTCS btcs.csv BTCS cn.csv CN

pip install pandas matplotlib numpy

# архитектура 
heat-equation/
│── cmd/
│   └── heat/          
│       └── main.go
│
│── internal/
│   ├── solver/        # численные методы (FTCS, BTCS, CN)
│   │   └── solver.go
│   │
│   ├── mathutils/     # аналитическое решение, начальные условия, ошибки
│   │   └── mathutils.go
│   │
│   ├── io/            # сохранение/загрузка данных (CSV)
│   │   └── io.go
│   │
│   └── config/        # параметры модели (структура Params)
│       └── config.go
