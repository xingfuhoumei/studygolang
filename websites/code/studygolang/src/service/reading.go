// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"strconv"
	"strings"

	"logger"
	"model"
)

// 获取晨读列表（分页）
func FindReadings(lastId, limit string, rtype int) []*model.MorningReading {
	reading := model.NewMorningReading()

	cond := "rtype=?"
	args := make([]interface{}, 0, 2)
	args = append(args, rtype)
	if lastId != "0" {
		cond += " AND id<?"
		args = append(args, lastId)
	}

	readingList, err := reading.Where(cond, args...).Order("id DESC").Limit(limit).
		FindAll()
	if err != nil {
		logger.Errorln("reading service FindReadings Error:", err)
		return nil
	}

	return readingList
}

// 【我要晨读】
func IReading(id string) string {
	_, err := strconv.Atoi(id)
	if err != nil {
		return "/readings"
	}

	reading := model.NewMorningReading()
	err = reading.Where("id=?", id).Find()

	if err != nil {
		logger.Errorln("reading service IReading error:", err)
		return "/readings"
	}

	if reading.Id == 0 {
		return "/readings"
	}

	go reading.Where("id=?", id).Increment("clicknum", 1)

	if reading.Inner == 0 {
		return "/wr?u=" + reading.Url
	}

	return "/articles/" + strconv.Itoa(reading.Inner)
}

// 获取晨读列表（分页）
func FindReadingByPage(conds map[string]string, curPage, limit int) ([]*model.MorningReading, int) {
	conditions := make([]string, 0, len(conds))
	for k, v := range conds {
		conditions = append(conditions, k+"="+v)
	}

	reading := model.NewMorningReading()

	limitStr := strconv.Itoa((curPage-1)*limit) + "," + strconv.Itoa(limit)
	readingList, err := reading.Where(strings.Join(conditions, " AND ")).Order("id DESC").Limit(limitStr).
		FindAll()
	if err != nil {
		logger.Errorln("reading service FindArticleByPage Error:", err)
		return nil, 0
	}

	total, err := reading.Count()
	if err != nil {
		logger.Errorln("reading service FindReadingByPage COUNT Error:", err)
		return nil, 0
	}

	return readingList, total
}
