import { Badge } from '@/components/atoms/Badge'
import { Text } from '@/components/atoms/Text'
import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

type SearchType = 'tag' | 'keyword'
type SearchProps = {
  type: SearchType
  key: string
}

export const SearchPage = () => {
  // TODO: 検索表示
  const [search, setSearch] = useState<SearchProps>()
  const [searchParams] = useSearchParams()

  useEffect(() => {
    if (searchParams.get('tag') !== null) {
      setSearch({ type: 'tag', key: searchParams.get('tag') ?? '' })
    }
    if (searchParams.get('keyword') !== null) {
      setSearch({
        type: 'keyword',
        key: searchParams.get('keyword') ?? '',
      })
    }
  }, [searchParams])

  return (
    <>
      <Text
        fontSize="large"
        color="gray"
        key={search?.key}
      >{`Searched by ${search?.type}: `}</Text>
      {search?.type === 'tag' ? (
        <Badge>{search.key}</Badge>
      ) : (
        <Badge backgroundColor="pink">{search?.key}</Badge>
      )}
    </>
  )
}
