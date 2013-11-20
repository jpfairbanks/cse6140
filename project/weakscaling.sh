file="weakscaling.txt"
if [ ! -r $file ]
 then
  for i in {1,2,4,8,16,64};
   do export GOMAXPROCS=$i;
    go test -test.bench=Weak -test.run=XXX; 
  done > $file
 else
  echo "cannot find " $file
fi
cat weakscaling.txt | grep ns | tr -s " " "\t" | cut -f3 |\
 awk 'BEGIN{poweroftwo=1/2} {poweroftwo*=2; print $1 /(poweroftwo * 1000) , "ms", poweroftwo}'
