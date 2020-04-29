
## 功能

    条件: 索引为name-yyyy.MM.dd, 时间字段为@timestamp

- 指定时间范围，字段等于特定值出现的次数超阈值告警  
   
   `count(A=xx;B=xx) > threshold in x minute`
- 指定时间范围，某个字段大于特定值出现的次数超阈值告警,可选支持排除条件 

  `count(A>xx) && !(B=yy;B=zz) > threshold in x minute`
- 指定时间范围，按某个字段聚合分组后，前topX出现次数超阈值告警,可选支持排除某组  
  
  `count(top(agg(A))) && !(A=xx) > threshold in  5minute`


## elastic数据迁移

[eslasticdump](https://github.com/taskrabbit/elasticsearch-dump?utm_source=dbweekly&utm_medium=email)

```bash
# install npm
cd /opt/
wget https://nodejs.org/dist/v10.16.0/node-v10.16.0-linux-x64.tar.xz
xz -d node-v10.16.0-linux-x64.tar.xz
tar -xvf node-v10.16.0-linux-x64.tar
ln -s /opt/nodejs/node-v10.16.0-linux-x64/bin/node /usr/local/bin/node
ln -s /opt/nodejs/node-v10.16.0-linux-x64/bin/npm /usr/local/bin/npm
cp ~/.bash_profile /home/bash_profile
echo 'export PATH=$PATH:/opt/node-v10.16.0-linux-x64/bin' >> ~/.bash_profile
source ~/.bash_profile

# global  install elasticdump 
npm install elasticdump -g

# test
elasticdump --help

# migration
elasticdump --input=http://elastic:elastic@10.11.40.66:9200/gitlab-2020.04.28 --output=gitlab_mapping.json --type=mapping
elasticdump --input=http://elastic:elastic@10.11.40. 66:9200/gitlab-2020.04.28 
--output=gitlab_data.json 
--type=data

elasticdump --input=http://elastic:elastic@10.11.40.66:9200/gitlab-2020.04.28 --output=http://10.51.1.31:9201/gitlab-2020.04.28 --type=mapping
elasticdump --input=http://elastic:elastic@10.11.40.66:9200/gitlab-2020.04.28 --output=http://10.51.1.31:9201/gitlab-2020.04.28 --type=data

```