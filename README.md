
## 功能

    条件: 索引为name-yyyy.MM.dd, 时间字段为@timestamp

- 指定时间范围，字段等于特定值出现的次数超阈值告警  
   
   `count(A=xx;B=xx) > threshold in x minute`
- 指定时间范围，某个字段大于特定值出现的次数超阈值告警,可选支持排除条件 

  `count(A>xx) && !(B=yy;B=zz) > threshold in x minute`
- 指定时间范围，按某个字段聚合分组后，前topX出现次数超阈值告警,可选支持排除某组  
  
  `count(top(agg(A))) && !(A=xx) > threshold in  5minute`

