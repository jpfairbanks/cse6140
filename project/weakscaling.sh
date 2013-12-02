file="weakscaling.txt"
for i in {1,2,4,8,10,16,20,24,32,40,48,64};
    do export GOMAXPROCS=$i;
    go test -test.bench=Weak -test.run=XXX -test.benchtime=3s;
done | tee $file
cat $file | grep ns | tr -s " " "\t" | cut -f3 | tee $file.post
#awk 'BEGIN{poweroftwo=1/2} {poweroftwo*=2; print $1 /(poweroftwo * 1000) , "ms", poweroftwo}' | tee $file.post
