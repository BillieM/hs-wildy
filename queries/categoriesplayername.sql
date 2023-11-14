-- SQLite
    select c.(*), p.name 
      from categories c 
 left join players p 
        on c.player_id = p.id;