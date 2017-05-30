# xtable

xtable is a simple sql query compostition tool.  Like Apache Pig, or old good
Ingres QUEL, it allow user to build up a complex query step by step.  Many users
prefer such an iterative approach and it is much easier/more natural for tool
writers.

xtable build next level query by doing a very simple "macro substitue".   
```
#x# -> where x is a number.  This is a table alias, will be substited by the x-th input table.
#x.y# -> where x is a number, y is eaither a number or a column name.  This will be substitued
         by x-th input table, y-th column (or by column name).
## -> # escape, will be replace by a single # char.
```

xtable is unique compared to other dataframe library in that
1. Access to full SQL syntax.  Including advanced feature like grouping, limit, sample, olap, etc.
2. Each step the query is syntax/type checked so that user can discover/fix errors early.
3. Completely lazy.  No execution at all until Execute or Render is called.  Predicate added in
   later composition steps can be pushed down and/or cause execute plan change.
   
xtable perform syntax/type checking by parsing explain output and heavily levarages WITH clause. 
While the database connection object is called Deepgreen, it is completely compatible with Greenplum.
A PostgreSQL port should not be too hard as well (TODO).
