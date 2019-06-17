package infrastructure

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main/domain"
)

type TabelRepository interface {
	GetChildTabel(ctx context.Context, childID primitive.ObjectID, dateFrom int64, dateTo int64) ([]*domain.TabelRecord, error)
	//getTabelWithPrices(ctx context.Context, childID primitive.ObjectID, dateFrom int64, dateTo int64)
	AddDay(ctx context.Context, date int64, childID primitive.ObjectID, status bool) (primitive.ObjectID, error)
}

type tabelRepository struct{}
func NewtabelRepository() TabelRepository {
	return &tabelRepository{}
}

//Python
//def get_child_tabel(self, child_id, date_from, date_to):
//self.cur.execute(r"SELECT * from mydb.tabel where date BETWEEN date(%s) AND date(%s) and childid = %s order by date", (date_from, date_to, child_id, ))
//rows = self.cur.fetchall()
//return rows
func (tr *tabelRepository) GetChildTabel(ctx context.Context, childID primitive.ObjectID, dateFrom int64, dateTo int64) ([]*domain.TabelRecord, error){

	database, err := InitDb(ctx)
	if err != nil {
		return nil, err
	}

	collection := database.Collection(TabelCollectionName)
	filter :=  bson.D{{"childID", childID},{"Date",bson.D{{"$gte",dateFrom}, {"$lte",dateTo}}}}

	tabelRecords, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var TabelRecords []*domain.TabelRecord

	for tabelRecords.Next(ctx) {
		var elem domain.TabelRecord
		err := tabelRecords.Decode(&elem)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		TabelRecords = append(TabelRecords, &elem)
	}

	err = database.Client().Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	return TabelRecords, nil
}

//Python
//def get_tabel_with_prices(self, child_id, date_from, date_to):
//
//query = r"""SELECT  tabel.`date`, tabel.childid, tabel.`value`, max(price.price_date), price.price, chil.`name`, chil.userid, chil.groupid
//FROM tabel
//LEFT JOIN mydb.children AS chil ON childid = chil.id
//LEFT JOIN price ON tabel.`date` >= price.price_date AND chil.groupid = price.groupid
//WHERE tabel.`value` = 1 AND tabel.childid = %s AND tabel.`date` BETWEEN date(%s) and date(%s)
//GROUP BY mydb.tabel.`date`, mydb.tabel.childid, mydb.chil.userid, mydb.chil.groupid"""
//
//self.cur.execute(query, (child_id, date_from, date_to, ))
//rows = self.cur.fetchall()
//return rows

//func (tr *tabelRepository) getTabelWithPrices(ctx context.Context, childID primitive.ObjectID, dateFrom int64, dateTo int64){
//
//
//}

//Python
//def add_day(self, date, child_id, status):
//self.cur.execute(r"SELECT * from mydb.tabel where date = date(%s) and childid = %s", (date, child_id, ))
//rows = self.cur.fetchall()
//if len(rows):
//self.cur.execute(r"UPDATE mydb.tabel set value = %s where date = date(%s) and childid = %s ", (status, date, child_id,))
//else:
//self.cur.execute(r"INSERT into mydb.tabel VALUES (NULL, date(%s), %s, %s)", (date, status, child_id,))
//self.conn.commit()

func (tr *tabelRepository) AddDay(ctx context.Context, date int64, childID primitive.ObjectID, status bool) (primitive.ObjectID, error) {

	var emptyVal primitive.ObjectID
	database, err := InitDb(ctx)
	if err != nil {
		return emptyVal, err
	}

	collection := database.Collection(TabelCollectionName)

	tabelRecord := domain.NewTabelRecord()
	tabelRecord.ID = primitive.NewObjectID()
	tabelRecord.Date = date
	tabelRecord.ChildID = childID
	tabelRecord.Value = status

	insertResult, err := collection.InsertOne(ctx, tabelRecord)
	if err != nil {
		return emptyVal, err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	err = database.Client().Disconnect(ctx)
	if err != nil {
		return emptyVal, err
	}

	return insertResult.InsertedID.(primitive.ObjectID), nil

}