import { Button } from '@/components/atoms/Button'
import Dropzone from '@/components/atoms/Dropzone'
import { Input } from '@/components/atoms/Input'
import { Text } from '@/components/atoms/Text'
import TextArea from '@/components/atoms/TextArea'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import TagForm from '@/components/molecules/TagForm'
import { useState } from 'react'
import { Controller, useForm } from 'react-hook-form'
import styled from 'styled-components'

export type BlogFormData = {
  title: string
  description: string
  content: string
  authorId: number
  isPublic: boolean
  thumbnailImageFileName?: string
  tags: string[]
}

type BlogFormProps = {
  data?: BlogFormData
  onSubmit?: (data: BlogFormData) => void
}

const PreviewImageWrapper = styled.div`
  width: 100%;
  height: 150px;
  > img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
`

export const BlogForm = (props: BlogFormProps) => {
  // TODO: isPublic
  // TODO: authorId
  const { data, onSubmit } = props

  const {
    control,
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<BlogFormData>({
    defaultValues: data,
  })

  const [imageFiles, setImageFiles] = useState<File[]>([])

  const handleOnSubmit = async (data: BlogFormData) => {
    data.isPublic = true
    data.authorId = 1
    console.log('submit')
    console.log(data)
    onSubmit && onSubmit(data)
  }

  return (
    <form>
      <Box>
        <Text as="label" variant="medium">
          Title
        </Text>
        <Box marginTop={1}>
          <Input
            {...register('title', { required: 'タイトルは必須です。' })}
            name="title"
            placeholder="Title"
            hasError={!!errors.title}
          />
          {errors.title && (
            <Text as="label" variant="small" color="danger">
              {errors.title.message}
            </Text>
          )}
        </Box>
      </Box>
      <Box marginTop={3}>
        <Text as="label" variant="medium">
          Description
        </Text>
        <Box marginTop={1}>
          <Input
            {...register('description', { required: '概要は必須です。' })}
            name="description"
            placeholder="Description"
            hasError={!!errors.description}
          />
          {errors.description && (
            <Text as="label" variant="small" color="danger">
              {errors.description.message}
            </Text>
          )}
        </Box>
      </Box>
      <Flex
        marginTop={3}
        flexDirection="row"
        alignItems="start"
        justifyContent="space-between"
      >
        <Box width="40%">
          <Text as="label" variant="medium">
            Thumbnail
          </Text>
          <Box marginTop={1}>
            <Controller
              control={control}
              name="thumbnailImageFileName"
              render={({ field: { value, onChange } }) => (
                <>
                  <Dropzone
                    value={imageFiles}
                    onChange={(files) => {
                      // TODO: 画像アップロード
                      if (files.length > 1) {
                        control.setError('thumbnailImageFileName', {
                          message: 'サムネイルは1つまでです。',
                        })
                        return
                      }
                      const url = URL.createObjectURL(files[0])
                      setImageFiles([files[0]])
                      onChange(url)
                    }}
                  >
                    {imageFiles.length > 0 && (
                      <Box>
                        <Box backgroundColor="primary" padding="3px">
                          <Text color="white">{imageFiles[0].name}</Text>
                        </Box>
                        <PreviewImageWrapper>
                          <img src={value} />
                        </PreviewImageWrapper>
                      </Box>
                    )}
                  </Dropzone>
                  {errors.thumbnailImageFileName && (
                    <Text as="label" variant="small" color="danger">
                      {errors.thumbnailImageFileName.message}
                    </Text>
                  )}
                </>
              )}
            />
          </Box>
        </Box>
        <Box width="55%">
          <Text as="label" variant="medium">
            Tags
          </Text>
          <Box marginTop={1}>
            <Controller
              control={control}
              name="tags"
              rules={{
                validate: (value) => {
                  return (
                    (value && value.length <= 3) ||
                    '選択できるタグは3つまでです。'
                  )
                },
              }}
              render={({ field: { onChange, value } }) => (
                <TagForm
                  placeholder="Tags"
                  value={value}
                  onKeyDown={(tags: string[]) => onChange(tags)}
                />
              )}
            />
            {errors.tags && (
              <Text as="label" variant="small" color="danger">
                {errors.tags.message}
              </Text>
            )}
          </Box>
        </Box>
      </Flex>
      <Box marginTop={3}>
        <Text as="label" variant="medium">
          Content
        </Text>
        <Box marginTop={1}>
          <Controller
            control={control}
            name="content"
            rules={{ validate: (value) => !!value || '本文は必須です。' }}
            render={({ field: { onChange, value } }) => (
              <TextArea minRows={10} value={value} onChange={onChange} />
            )}
          />
          {errors.content && (
            <Text as="label" variant="small" color="danger">
              {errors.content.message}
            </Text>
          )}
        </Box>
      </Box>
      <Flex justifyContent="flex-end" marginTop={2}>
        <Button
          variant="primary"
          type="button"
          onClick={handleSubmit(handleOnSubmit)}
        >
          Post
        </Button>
      </Flex>
    </form>
  )
}
